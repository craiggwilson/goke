package task

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-isatty"

	"github.com/craiggwilson/goke/task/internal"
)

const trueString = "true"

// Run orders the tasks be dependencies to build an execution plan and then executes each required task.
func Run(registry *Registry, arguments []string) error {
	opts, err := parseArgs(arguments)
	if err != nil {
		return err
	}

	ui := newTUI(opts.color)

	if opts.help {
		return printHelp(ui, registry)
	}

	tasksToRun, err := sortTasksToRun(registry.Tasks(), opts.taskNames)
	if err != nil {
		return err
	}

	if len(tasksToRun) == 0 {
		return printHelp(ui, registry)
	}

	writer := internal.NewPrefixWriter(&syncWriter{Writer: os.Stdout})
	prefix := []byte("       | ")

	unusedArgs := getUnusedArgs(tasksToRun, opts.args)
	if len(unusedArgs) > 0 {
		for _, unusedArg := range unusedArgs {
			_, _ = fmt.Fprintln(writer, ui.Error("WARNING"), "unused argument", unusedArg)
		}
		if registry.shouldErrorOnUnusedArgs {
			return fmt.Errorf("unused args")
		}
	}

	totalStartTime := time.Now()

	var failedTasks []string
	var finallyTaskNames []string
	for _, t := range tasksToRun {
		executor := t.Executor()
		if executor == nil {
			// this task is just an aggregate task
			finallyTaskNames = append(t.Finally(), finallyTaskNames...)
			continue
		}

		taskArgs, err := argsForTask(t, opts.args)
		if err != nil {
			return err
		}

		finallyTaskNames = append(t.Finally(), finallyTaskNames...)

		ctx := NewContext(context.Background(), writer, taskArgs)
		ctx.UI = ui
		ctx.Verbose = opts.verbose

		ctx.Logln(ui.Info("START"), " |", ui.Highlight(t.Name()))
		writer.SetPrefix(prefix)

		startTime := time.Now()
		err = executor(ctx)
		finishedTime := time.Now()

		writer.SetPrefix(nil)
		if err != nil {
			failedTasks = append(failedTasks, t.Name())
			ctx.Logln(ui.Error("FAIL"), "  |", ui.Highlight(t.Name()), "in", finishedTime.Sub(startTime).String())
			writer.SetPrefix(prefix)
			ctx.Logln(ui.Highlight(err.Error()))
			writer.SetPrefix(nil)
			if !t.ContinueOnError() {
				break
			}
		} else {
			ctx.Logln(ui.Success("FINISH"), "|", ui.Highlight(t.Name()), "in", finishedTime.Sub(startTime).String())
		}
	}

	if finallyTasks, err := sortTasksToRun(registry.Tasks(), finallyTaskNames); err == nil && len(finallyTasks) > 0 {
		ctx := NewContext(context.Background(), writer, nil)
		ctx.UI = ui
		ctx.Verbose = opts.verbose
		ctx.Logln(ui.Info("START"), " |", ui.Highlight("finalizing tasks"))
		writer.SetPrefix(prefix)
		startTime := time.Now()
		for _, task := range finallyTasks {
			if task.Executor() != nil {
				if err := task.Executor()(ctx); err != nil {
					writer.SetPrefix(nil)
					ctx.Logln(ui.Warning("WARN"), "  |", ui.Highlight(task.Name()), "failed:", err.Error())
					writer.SetPrefix(prefix)
				} else {
					ctx.Logln(ui.Highlight(task.Name()), "finished")
				}
			}
		}
		writer.SetPrefix(nil)
		ctx.Logln(ui.Success("FINISH"), "|", ui.Highlight("finalizing tasks"), "in", time.Since(startTime).String())
	} else if err != nil {
		// should not happen since Finally() tasks are validated when building the primary task list
		fmt.Fprintln(writer, ui.Error("WARNING"), "Building finalizing task list failed:", err.Error())
	}

	totalDuration := time.Since(totalStartTime)

	if len(failedTasks) > 0 {
		return fmt.Errorf("task(s) %s failed", failedTasks)
	}

	_, _ = fmt.Fprintln(writer, "---------------")
	_, _ = fmt.Fprintln(writer, ui.Success(fmt.Sprint("Completed in ", totalDuration)))

	return nil
}

func argsForTask(task Task, args globalArgs) (map[string]string, error) {
	taskArgs := make(map[string]string)
	for _, da := range task.DeclaredArgs() {
		// first look up a specific one to the task
		v, ok := args.get(task.Name(), da.Name)
		if !ok {
			// try to find one in the global namespace
			v, ok = args.get("", da.Name)
		}

		if da.Validator != nil {
			if err := da.Validator(da.Name, v); err != nil {
				return nil, fmt.Errorf("failed to validate argument %q: %v", da.Name, err)
			}
		}

		if ok {
			taskArgs[da.Name] = v
		}
	}

	return taskArgs, nil
}

func parseArgs(arguments []string) (*runOptions, error) {
	var requiredTaskNames []string
	args := globalArgs{}
	for _, arg := range arguments {
		if arg[0] == '-' || arg[0] == '/' {
			taskName, argName, value := parseArg(arg)
			switch argName {
			case "h":
				argName = "help"
			case "v":
				argName = "verbose"
			}

			args.set(taskName, argName, value)
		} else {
			requiredTaskNames = append(requiredTaskNames, arg)
		}
	}

	verboseArg, _ := args.get("", "verbose")
	verbose := verboseArg == trueString
	helpArg, _ := args.get("", "help")
	help := helpArg == trueString

	color := isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
	if colorArg, ok := args.get("", "color"); ok && colorArg != trueString {
		color = false
	}

	return &runOptions{
		args:      args,
		verbose:   verbose,
		help:      help,
		color:     color,
		taskNames: requiredTaskNames,
	}, nil
}

func getUnusedArgs(tasks []Task, args globalArgs) []string {
	var used = make(map[string]map[string]bool)

	for ns, nsArgs := range args {
		used[ns] = make(map[string]bool, len(nsArgs))
		for arg := range nsArgs {
			used[ns][arg] = false
		}
	}

	for _, task := range tasks {
		for _, da := range task.DeclaredArgs() {
			if _, ok := args.get(task.Name(), da.Name); ok {
				used[task.Name()][da.Name] = true
			} else if _, ok = args.get("", da.Name); ok {
				used[""][da.Name] = true
			}
		}
	}

	// Check to make sure that everything is used.
	var unusedArgs []string
	for ns, args := range used {
		for arg, didUse := range args {
			if !didUse {
				if ns == "" {
					unusedArgs = append(unusedArgs, arg)
				} else {
					unusedArgs = append(unusedArgs, fmt.Sprintf("%s:%s", ns, arg))
				}
			}
		}
	}

	return unusedArgs
}

func parseArg(arg string) (string, string, string) {
	arg = strings.TrimLeftFunc(arg, func(r rune) bool {
		return r == '-' || r == '/'
	})
	parts := strings.SplitN(arg, "=", 2)
	ns, name := parseArgName(parts[0])
	if len(parts) == 1 {
		return ns, name, trueString
	}

	return ns, name, parts[1]
}

func parseArgName(name string) (string, string) {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) == 1 {
		return "", parts[0]
	}

	return parts[0], parts[1]
}

func printHelp(ui *TUI, registry *Registry) error {
	fs := flag.NewFlagSet("goke", flag.ContinueOnError)
	_ = fs.Bool("v", false, "generate verbose logs")
	usage(ui, fs, registry)
	return flag.ErrHelp
}

type runOptions struct {
	args      globalArgs
	verbose   bool
	help      bool
	color     bool
	taskNames []string
}

type globalArgs map[string]map[string]string

func (ga globalArgs) get(taskName, argName string) (string, bool) {
	if ta, ok := ga[taskName]; ok {
		if v, ok := ta[argName]; ok {
			return v, true
		}
	}
	return "", false
}

func (ga globalArgs) set(taskName, argName, value string) {
	ta, ok := ga[taskName]
	if !ok {
		ta = make(map[string]string)
		ga[taskName] = ta
	}

	ta[argName] = value
}

type syncWriter struct {
	io.Writer

	mu sync.Mutex
}

func (w *syncWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Writer.Write(p)
}
