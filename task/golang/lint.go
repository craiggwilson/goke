package golang

import (
	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

// LintOptions represents the arguments to the golint command.
type LintOptions struct {
	Paths         []string
	SetExitStatus bool
}

// Lint returns a function that runs golint.
func Lint(opts *LintOptions) task.Executor {
	args := []string{}
	if opts.SetExitStatus {
		args = append(args, "-set_exit_status")
	}
	args = append(args, opts.Paths...)
	return command.Command("golint", args...)
}
