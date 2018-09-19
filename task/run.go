package task

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

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

	prefix := []byte("      | ")

	totalStartTime := time.Now()

	for _, t := range tasksToRun {
		if t.Executor() == nil {
			// this task is just an aggregate task
			continue
		}
		ctx.Logln("START |", t.Name())
		writer.SetPrefix(prefix)
		startTime := time.Now()
		err := t.Executor()(ctx)
		finishedTime := time.Now()
		writer.SetPrefix(nil)
		if err != nil {
			ctx.Logln("FAIL  |", t.Name())
			writer.SetPrefix(prefix)
			ctx.Logln(err)
			return fmt.Errorf("task '%s' failed", t.Name())
		}
		ctx.Logf("FINISH| %s in %v\n", t.Name(), finishedTime.Sub(startTime))
	}

	totalDuration := time.Now().Sub(totalStartTime)

	ctx.Logln("---------------")
	ctx.Logln("Completed in ", totalDuration)

	return nil
}

func parseArgs(registry *Registry, arguments []string) (*runOptions, error) {

	requiredTaskNames := parseRequiredTaskNames(arguments)
	arguments = arguments[len(requiredTaskNames):]

	fs := flag.NewFlagSet("goke", flag.ContinueOnError)
	fs.Usage = func() {
		usage(fs, registry)
	}
	dryrun := fs.Bool("dryrun", false, "performs a dry run, executing each task with the dry-run flag")
	verbose := fs.Bool("v", false, "generate verbose logs")
	if err := fs.Parse(arguments); err != nil {
		return nil, err
	}

	extraArgs := fs.Args()

	return &runOptions{
		dryrun:    *dryrun,
		extraArgs: extraArgs,
		verbose:   *verbose,
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
