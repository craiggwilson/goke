package task

import (
	"context"
	"flag"
	"fmt"
	"os"
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
	ctx := &Context{
		Args:    opts.extraArgs,
		Context: context.Background(),
		DryRun:  opts.dryrun,
		Verbose: opts.verbose,
		w:       writer,
	}

	prefix := []byte("       | ")

	totalStartTime := time.Now()

	cBright := ansi.ColorFunc("white+bh")
	cInfo := ansi.ColorFunc("cyan+b")
	cFail := ansi.ColorFunc("red+b")
	cSuccess := ansi.ColorFunc("green+b")

	for _, t := range tasksToRun {
		if t.Executor() == nil {
			// this task is just an aggregate task
			continue
		}
		ctx.Logln(cInfo("START"), " |", cBright(t.Name()))
		writer.SetPrefix(prefix)
		startTime := time.Now()
		err := t.Executor()(ctx)
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

	ctx.Logln("---------------")
	ctx.Logln(cSuccess(fmt.Sprint("Completed in ", totalDuration)))

	return nil
}

func parseArgs(registry *Registry, arguments []string) (*runOptions, error) {
	dryrun := false
	verbose := false
	help := false
	var requiredTaskNames []string
	var extraArgs []string
	seenFlags := false
	for _, arg := range arguments {
		if arg[0] == '-' || arg[0] == '/' {
			seenFlags = true
			switch arg {
			case "-v", "--v", "/v":
				verbose = true
			case "-dryrun", "--dryrun", "/dryrun":
				dryrun = true
			case "-help", "--help", "/help", "-h", "--h", "/h":
				help = true
			default:
				extraArgs = append(extraArgs, arg)
			}
		} else {
			if !seenFlags {
				requiredTaskNames = append(requiredTaskNames, arg)
			} else {
				extraArgs = append(extraArgs, arg)
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
		extraArgs: extraArgs,
		verbose:   verbose,
		taskNames: requiredTaskNames,
	}, nil
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
	extraArgs []string
	verbose   bool
	taskNames []string
}
