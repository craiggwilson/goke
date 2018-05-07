package golang

import (
	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

// BuildOptions represents the arguments to the go build command.
type BuildOptions struct {
	Paths   []string
	Verbose bool
	Flags   []string
	Output  string
}

// Build returns a function that runs go build.
func Build(opts *BuildOptions) task.Executor {
	args := []string{}
	args = append(args, "build")
	if opts.Verbose {
		args = append(args, "-v")
	}
	if len(opts.Flags) > 0 {
		args = append(args, opts.Flags...)
	}
	if opts.Output != "" {
		args = append(args, []string{"-o", opts.Output}...)
	}
	args = append(args, opts.Paths...)
	return command.Command("go", args...)
}
