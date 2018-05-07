package task

import (
	"flag"
	"fmt"
	"strings"
)

func usage(fs *flag.FlagSet, registry *Registry) {
	fmt.Println("USAGE: [options ...] [tasks ...]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fs.PrintDefaults()
	fmt.Println()
	fmt.Println("TASKS:")
	for _, t := range registry.tasks {
		fmt.Print("  ", t.Name())
		if len(t.Dependencies()) > 0 {
			fmt.Printf(" -> (%s)", strings.Join(t.Dependencies(), ", "))
		}
		fmt.Println()
		if t.Description() != "" {
			fmt.Println("       ", t.Description())
		}
	}
}
