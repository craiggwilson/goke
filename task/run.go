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

		taskArgs := make(map[string]string)
		for _, da := range t.DeclaredArgs() {
			// first look up a specific one to the task
			v, ok := opts.args.get(t.Name(), da.Name)
			if !ok {
				// try to find one in the global namespace
				v, ok = opts.args.get("", da.Name)
			}

			if ok {
				taskArgs[da.Name] = v
			} else if da.Required {
				return fmt.Errorf("task %q has a required argument %q that was not provided", t.Name(), da.Name)
			}
		}

		ctx := &Context{
			Context: context.Background(),
			DryRun:  opts.dryrun,
			Verbose: opts.verbose,

			taskArgs: taskArgs,
			w:        writer,
		}

		ctx.Logln(cInfo("START"), " |", cBright(t.Name()))
		writer.SetPrefix(prefix)

		startTime := time.Now()
		err := executor(ctx)
		finishedTime := time.Now()

		writer.SetPrefix(nil)
		if err != nil {
			ctx.Logln(cFail("FAIL"), "  |", cBright(t.Name()))
			writer.SetPrefix(prefix)
			ctx.Logln(cBright(err.Error()))
			return fmt.Errorf("task '%s' failed", t.Name())
		}
		ctx.Logln(cSuccess("FINISH"), "|", cBright(fmt.Sprintf("%s in %v", t.Name(), finishedTime.Sub(startTime))))
	}

	totalDuration := time.Now().Sub(totalStartTime)

	fmt.Fprintln(writer, "---------------")
	fmt.Fprintln(writer, cSuccess(fmt.Sprint("Completed in ", totalDuration)))

	return nil
}

func parseArgs(registry *Registry, arguments []string) (*runOptions, error) {
	dryrun := false
	verbose := false
	help := false
	var requiredTaskNames []string
	args := globalArgs{}
	seenFlags := false
	for _, arg := range arguments {
		if arg[0] == '-' || arg[0] == '/' {
			seenFlags = true
			switch arg {
			case "--verbose", "-v", "/v":
				verbose = true
			case "--dryrun", "/dryrun":
				dryrun = true
			case "--help", "/help", "-h", "--h", "/h":
				help = true
			default:
				taskName, argName, value := parseExtraArg(arg)
				args.set(taskName, argName, value)
			}
		} else {
			if !seenFlags {
				requiredTaskNames = append(requiredTaskNames, arg)
			} else {
				taskName, argName, value := parseExtraArg(arg)
				args.set(taskName, argName, value)
			}
		}
	}

	if help {
		fs := flag.NewFlagSet("goke", flag.ContinueOnError)
		_ = fs.Bool("dryrun", false, "performs a dry run, executing each task with the dry-run flag")
		_ = fs.Bool("v", false, "generate verbose logs")
		usage(fs, registry)
		return nil, flag.ErrHelp
	}

	return &runOptions{
		dryrun:    dryrun,
		args:      args,
		verbose:   verbose,
		taskNames: requiredTaskNames,
	}, nil
}

func parseExtraArg(arg string) (string, string, string) {
	arg = strings.TrimLeftFunc(arg, func(r rune) bool {
		return r == '-' || r == '/'
	})
	parts := strings.SplitN(arg, "=", 2)
	ns, name := parseExtraArgName(parts[0])
	if len(parts) == 1 {
		return ns, name, "true"
	}

	return ns, name, parts[1]
}

func parseExtraArgName(name string) (string, string) {
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
	dryrun    bool
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
