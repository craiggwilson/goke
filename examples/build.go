package main

import (
	"fmt"
	"os"
	"time"

	"github.com/craiggwilson/goke/task"
	"github.com/craiggwilson/goke/task/golang"
)

var registry = task.NewRegistry()

func init() {
	registry.Declare("First").
		Do(func(ctx *task.Context) error {
			ctx.Logln("HERE")
			time.Sleep(2 * time.Second)
			return nil
		})

	registry.Declare("Second").
		DependsOn("First").
		Do(func(ctx *task.Context) error {
			ctx.Logln("THERE\n another line\nawesome")
			return nil
		})

	registry.Declare("Build").
		Do(golang.Build(&golang.GoBuild{
			Packages: []string{"./examples/simple.go"},
		}))
}

func main() {
	err := task.Run(registry, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
