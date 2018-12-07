package task

import (
	"flag"
	"fmt"
)

func usage(fs *flag.FlagSet, registry *Registry) {
	fmt.Println("USAGE: [tasks ...] [options ...]")
	fmt.Println()
	fmt.Println("TASKS:")
	for _, t := range registry.tasks {
		if t.Hidden() {
			continue
		}
		fmt.Print("  ", t.Name())
		args := t.DeclaredArgs()
		if len(args) > 0 {
			fmt.Print("(")
			for i, a := range args {
				fmt.Print(a.Name)
				if i < len(args)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Print(")")
		}
		if len(t.Dependencies()) > 0 {
			fmt.Print(" -> ", t.Dependencies())
		}
		fmt.Println()
		if t.Description() != "" {
			fmt.Println("       ", t.Description())
		}
	}
	fmt.Println()
	fmt.Println("OPTIONS:")
	fs.PrintDefaults()
}
