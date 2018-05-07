package build

import (
	"os"

	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/golang"
)

func Build(ctx *task.Context) error {
	return golang.Build(&golang.BuildOptions{
		Paths:   []string{mainFile},
		Output:  buildOutputFile,
		Verbose: ctx.Verbose,
	})(ctx)
}

func Clean(ctx *task.Context) error {
	_ = os.Remove(buildOutputFile)
	return nil
}

func Fmt(ctx *task.Context) error {
	return golang.Fmt(&golang.FmtOptions{
		Paths:     packages,
		AllErrors: ctx.Verbose,
		Simplify:  true,
	})(ctx)
}

func Lint(ctx *task.Context) error {
	return golang.Lint(&golang.LintOptions{
		Paths: packages,
	})(ctx)
}

func Vet(ctx *task.Context) error {
	return golang.Vet(&golang.VetOptions{
		Paths:   packages,
		Verbose: ctx.Verbose,
	})(ctx)
}
