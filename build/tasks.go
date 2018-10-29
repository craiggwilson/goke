package build

import (
	"os"

	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/command"
)

func Build(ctx *task.Context) error {
	args := []string{"build", "-o", buildOutputFile}
	if ctx.Verbose {
		args = append(args, "-v")
	}

	args = append(args, mainFile)
	return command.Command("go", args...)(ctx)
}

func Clean(ctx *task.Context) error {
	_ = os.Remove(buildOutputFile)
	return nil
}

func Fmt(ctx *task.Context) error {
	args := []string{"-s", "-l"}
	if ctx.Verbose {
		args = append(args, "-e")
	}

	args = append(args, mainFile)
	return command.Command("gofmt", args...)(ctx)
}

func Lint(ctx *task.Context) error {
	args := []string{"-set_exit_status"}
	args = append(args, packages...)
	return command.Command("golint", args...)(ctx)
}

func Test(ctx *task.Context) error {
	args := []string{"test"}
	if ctx.Verbose {
		args = append(args, "-v")
	}
	args = append(args, packages...)
	return command.Command("go", args...)(ctx)
}
