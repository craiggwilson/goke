package golang

import (
	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

// GoBuild represents the arguments to the go build command.
type GoBuild struct {
	Packages []string
}

// Build returns a function that runs go build.
func Build(build *GoBuild) func(*task.Context) error {
	args := []string{"build"}
	args = append(args, build.Packages...)
	return command.Command("go", args...)
}
