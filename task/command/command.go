package command

import (
	"os/exec"
	"strings"

	"github.com/craiggwilson/goke/task"
)

// Command wraps exec.Command in a task executor.
func Command(name string, args ...string) task.Executor {
	return Executor(exec.Command(name, args...))
}

// Executor creates a task.Executor from the command.
func Executor(cmd *exec.Cmd) task.Executor {
	return func(ctx *task.Context) error {
		ctx.Logf("exec: '%s %s'\n", cmd.Path, strings.Join(cmd.Args[1:], " "))

		if !ctx.DryRun {
			cmd.Stdout = ctx
			cmd.Stderr = ctx
			cmd.Start()
			return cmd.Wait()
		}

		return nil
	}
}
