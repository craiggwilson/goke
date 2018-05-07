package golang

import (
	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

// VetOptions represents the arguments to the go vet command.
type VetOptions struct {
	Paths   []string
	Verbose bool
}

// Vet returns a function that runs go vet.
func Vet(opts *VetOptions) task.Executor {
	args := []string{"vet"}
	if opts.Verbose {
		args = append(args, "-x")
	}
	args = append(args, opts.Paths...)
	return command.Command("go", args...)
}
