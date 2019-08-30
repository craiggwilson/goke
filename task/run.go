package task

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mgutz/ansi"

	"github.com/craiggwilson/goke/task/internal"
)

// Run orders the tasks be dependencies to build an execution plan and then executes each required task.
func Run(registry *Registry, arguments []string) error {
	opts, err := parseArgs(registry, arguments)
	if err != nil {
		return err
	}

	tasksToRun, err := sort(registry.tasks, opts.taskNames)
	if err != nil {
		return err
	}

	if len(tasksToRun) == 0 {
		_, err = parseArgs(registry, append(arguments, "-h"))
		return err
	}

	writer := internal.NewPrefixWriter(os.Stdout)
	prefix := []byte("       | ")

	totalStartTime := time.Now()

	cBright := ansi.ColorFunc("white+bh")
	cInfo := ansi.ColorFunc("cyan+b")
	cFail := ansi.ColorFunc("red+b")
	cSuccess := ansi.ColorFunc("green+b")

	for _, t := range tasksToRun {
		executor := t.Executor()
		if executor == nil {
			// this task is just an aggregate task
			continue
		}

		taskArgs, err := argsForTask(t, opts.args)
		if err != nil {
			return err
		}

		ctx := NewContext(context.Background(), writer, taskArgs)
		ctx.Verbose = opts.verbose

		ctx.Logln(cInfo("START"), " |", cBright(t.Name()))
		writer.SetPrefix(prefix)

		startTime := time.Now()
		err = executor(ctx)
		finishedTime := time.Now()

		writer.SetPrefix(nil)
		if err != nil {
			ctx.Logln(cFail("FAIL"), "  |", cBright(t.Name()))
			writer.SetPrefix(prefix)
			ctx.Logln(cBright(err.Error()))
			return fmt.Errorf("task %q failed", t.Name())
		}
		ctx.Logln(cSuccess("FINISH"), "|", cBright(fmt.Sprintf("%s in %v", t.Name(), finishedTime.Sub(startTime))))
	}

	totalDuration := time.Now().Sub(totalStartTime)

	fmt.Fprintln(writer, "---------------")
	fmt.Fprintln(writer, cSuccess(fmt.Sprint("Completed in ", totalDuration)))

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

func parseArgs(registry *Registry, arguments []string) (*runOptions, error) {
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
	verbose := verboseArg == "true"
	helpArg, _ := args.get("", "help")
	if helpArg == "true" {
		fs := flag.NewFlagSet("goke", flag.ContinueOnError)
		_ = fs.Bool("v", false, "generate verbose logs")
		usage(fs, registry)
		return nil, flag.ErrHelp
	}

	return &runOptions{
		args:      args,
		verbose:   verbose,
		taskNames: requiredTaskNames,
	}, nil
}

func parseArg(arg string) (string, string, string) {
	arg = strings.TrimLeftFunc(arg, func(r rune) bool {
		return r == '-' || r == '/'
	})
	parts := strings.SplitN(arg, "=", 2)
	ns, name := parseArgName(parts[0])
	if len(parts) == 1 {
		return ns, name, "true"
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

func parseRequiredTaskNames(arguments []string) []string {
	var requiredTaskNames []string
	for _, arg := range arguments {
		if arg[0] == '-' || arg[0] == '/' {
			break
		}

		requiredTaskNames = append(requiredTaskNames, arg)
	}

	return requiredTaskNames
}

type runOptions struct {
	args      globalArgs
	verbose   bool
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
