package golang

import (
	"strconv"
	"time"

	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

// TestOptions represents the arguments to the go test command.
type TestOptions struct {
	Failfast bool
	List     string
	Parallel int
	Paths    []string
	Timeout  time.Duration
	Verbose  bool
	VetList  string
}

// Test returns a function that runs go test.
func Test(opts *TestOptions) task.Executor {
	args := []string{"test"}
	if opts.Failfast {
		args = append(args, "-failfast")
	}
	if opts.Parallel > 0 {
		args = append(args, "-parallel", strconv.Itoa(opts.Parallel))
	}
	if opts.Timeout > 0 {
		args = append(args, "-timeout", opts.Timeout.String())
	}
	if opts.Verbose {
		args = append(args, "-v")
	}

	args = append(args, opts.Paths...)
	return command.Command("go", args...)
}
