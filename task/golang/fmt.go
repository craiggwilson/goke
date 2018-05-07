package golang

import (
	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

// FmtOptions represents the arguments to the gofmt command.
type FmtOptions struct {
	Paths     []string
	AllErrors bool
	Simplify  bool
}

// Fmt returns a function that runs gofmt.
func Fmt(opts *FmtOptions) task.Executor {
	args := []string{"-w"}
	if opts.AllErrors {
		args = append(args, "-e")
	}
	if opts.Simplify {
		args = append(args, "-s")
	}
	args = append(args, opts.Paths...)
	return command.Command("gofmt", args...)
}
