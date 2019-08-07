package build

import (
	"os"

	"github.com/craiggwilson/goke/pkg/sh"
	"github.com/craiggwilson/goke/task"
)

func Build(ctx *task.Context) error {
	args := []string{"build", "-o", buildOutputFile}
	if ctx.Verbose {
		args = append(args, "-v")
	}

	args = append(args, mainFile)
	return sh.Run(ctx, "go", args...)
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
	return sh.Run(ctx, "gofmt", args...)
}

func Lint(ctx *task.Context) error {
	args := []string{"-set_exit_status"}
	args = append(args, packages...)
	return sh.Run(ctx, "golint", args...)
}

func Test(ctx *task.Context) error {
	args := []string{"test"}
	if ctx.Verbose {
		args = append(args, "-v")
	}
	args = append(args, packages...)
	return sh.Run(ctx, "go", args...)
}
