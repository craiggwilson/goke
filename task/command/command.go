package command

import (
	"os/exec"
	"strings"

	"github.com/craiggwilson/goke/task"
)

// Command wraps exec.Command in a task executor.
func Command(name string, args ...string) func(*task.Context) error {
	return func(ctx *task.Context) error {
		cmd := exec.Command(name, args...)

		ctx.Logf("exec: '%s %s'\n", cmd.Path, strings.Join(cmd.Args[1:], " "))

		if !ctx.DryRun {
			cmd.Stdout = ctx.Writer()
			cmd.Stderr = ctx.Writer()
			return cmd.Run()
		}

		return nil
	}
}
