package task

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Run orders the tasks be dependencies to build an execution plan and then executes each required task.
func Run(registry *Registry, arguments []string) error {

	fs := flag.NewFlagSet("goke", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Println("Usage: [options ...] [tasks ...]")
		fs.PrintDefaults()
	}
	dryrun := fs.Bool("dryrun", false, "performs a dry run, executing each task with the dry-run flag")
	list := fs.Bool("list", false, "lists all the configured tasks")

	if err := fs.Parse(arguments[1:]); err != nil {
		return err
	}

	if *list {

	}

	var requiredTaskNames []string
	for i := 0; i < fs.NArg(); i++ {
		arg := fs.Arg(i)
		if arg[0] != '-' && arg[0] != '/' {
			requiredTaskNames = append(requiredTaskNames, arg)
		} else {
			break
		}
	}

	tasksToRun, err := orderTasks(registry.tasks, requiredTaskNames)
	if err != nil {
		return err
	}

	writer := &indentWriter{
		w:  os.Stdout,
		nl: true,
	}

	ctx := &Context{
		DryRun: *dryrun,
		w:      writer,
	}

	prefix := []byte("      | ")

	totalStartTime := time.Now()

	for _, t := range tasksToRun {
		ctx.Logln("START |", t.Name())
		writer.prefix = prefix
		startTime := time.Now()
		err := t.Execute(ctx)
		finishedTime := time.Now()
		writer.prefix = nil
		if err != nil {
			ctx.Logln("FAIL  |", t.Name())
			writer.prefix = prefix
			ctx.Logln(err)
			return err
		}
		ctx.Logf("FINISH| %s in %v\n", t.Name(), finishedTime.Sub(startTime))
	}

	totalDuration := time.Now().Sub(totalStartTime)

	ctx.Logln("---------------")
	ctx.Logln("Completed in ", totalDuration)

	return nil
}

type indentWriter struct {
	w      io.Writer
	prefix []byte
	nl     bool
}

func (iw *indentWriter) Write(p []byte) (n int, err error) {
	for _, c := range p {
		if iw.nl {
			_, err = iw.w.Write(iw.prefix)
			if err != nil {
				return n, err
			}
			iw.nl = false
			n += len(iw.prefix)
		}

		_, err = iw.w.Write([]byte{c})
		if err != nil {
			return n, err
		}

		n++
		iw.nl = c == '\n'
	}

	return n, nil
}

func orderTasks(allTasks []Task, requiredTaskNames []string) ([]Task, error) {
	graph, err := buildGraph(allTasks, requiredTaskNames)
	if err != nil {
		return nil, err
	}

	result, err := toposort(graph)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func buildGraph(allTasks []Task, requiredTaskNames []string) ([]*graphNode, error) {
	allTasksMap := make(map[string]Task)
	for _, t := range allTasks {
		allTasksMap[strings.ToLower(t.Name())] = t
	}

	var g []*graphNode
	seenTasks := make(map[string]struct{})
	for len(requiredTaskNames) > 0 {
		taskName := requiredTaskNames[0]
		requiredTaskNames = requiredTaskNames[1:]

		task, ok := allTasksMap[strings.ToLower(taskName)]
		if !ok {
			return nil, fmt.Errorf("unknown task '%s'", taskName)
		}

		if _, ok := seenTasks[task.Name()]; !ok {
			seenTasks[task.Name()] = struct{}{}
			g = append(g, &graphNode{task: task, edges: task.Dependencies()})
			requiredTaskNames = append(requiredTaskNames, task.Dependencies()...)
		}
	}

	return g, nil
}

type graphNode struct {
	task  Task
	edges []string
}

func toposort(g []*graphNode) ([]Task, error) {
	var queue []*graphNode
	for _, n := range g {
		if len(n.edges) == 0 {
			queue = append(queue, n)
		}
	}

	var sorted []Task
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		sorted = append(sorted, n.task)
		for _, m := range g {
			for i := range m.edges {
				if m.edges[i] == n.task.Name() {
					m.edges = append(m.edges[:i], m.edges[i+1:]...)
					if len(m.edges) == 0 {
						queue = append(queue, m)
					}
					break
				}
			}
		}
	}

	for _, n := range g {
		if len(n.edges) > 0 {
			return nil, fmt.Errorf("a cycle exists")
		}
	}

	return sorted, nil
}
