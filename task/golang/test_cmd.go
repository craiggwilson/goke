package golang

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

// TestOptions represents the arguments to the go test command.
type TestOptions struct {
	Count    int
	Failfast bool
	List     string
	Parallel int
	Paths    []string
	Tags     []string
	Timeout  time.Duration
	Verbose  bool
	VetList  string
}

// Test returns a function that runs go test.
func Test(opts *TestOptions) task.Executor {
	args := []string{"test"}
	if opts.Count > 0 {
		args = append(args, fmt.Sprintf("-count=%d", opts.Count))
	}
	if opts.Failfast {
		args = append(args, "-failfast")
	}
	if opts.Parallel > 0 {
		args = append(args, "-parallel", strconv.Itoa(opts.Parallel))
	}
	if len(opts.Tags) > 0 {
		args = append(args, fmt.Sprintf("-tags=%s", strings.Join(opts.Tags, ",")))
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
