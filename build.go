package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/craiggwilson/goke/build"
	"github.com/craiggwilson/goke/task"
)

var registry = task.NewRegistry()

func init() {
	registry.Declare("build").Description("build the goke build script").DependsOn("clean", "sa").Do(build.Build)
	registry.Declare("clean").Description("cleans up the artifacts").Do(build.Clean)
	registry.Declare("lint").Description("lint the packages").Do(build.Lint)
	registry.Declare("vet").Description("vets the packages").Do(build.Vet)
	registry.Declare("fmt").Description("formats the packages").Do(build.Fmt)
	registry.Declare("sa").Description("performs a static analysis").DependsOn("fmt", "lint", "vet")
}

func main() {
	err := task.Run(registry, os.Args)
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(1)
		}
		fmt.Println(err)
		os.Exit(2)
	}
}
