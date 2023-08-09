package task

import (
	"fmt"
	"strings"
)

type state = uint8

const (
	unvalidated state = iota
	validating
	valid
)

func sortTasksToRun(allTasks []Task, requiredTaskNames []string) ([]Task, error) {
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
	finallyTaskStates := make(map[string]state)
	for len(requiredTaskNames) > 0 {
		taskName := requiredTaskNames[0]
		requiredTaskNames = requiredTaskNames[1:]

		task, ok := allTasksMap[strings.ToLower(taskName)]
		if !ok {
			return nil, fmt.Errorf("unknown task '%s'", taskName)
		}

		if _, ok := seenTasks[task.Name()]; !ok {
			seenTasks[task.Name()] = struct{}{}
			if err := validateFinallyTasks(allTasksMap, finallyTaskStates, task.Finally()); err != nil {
				return nil, err
			}
			// toposort can
			g = append(g, &graphNode{task: task, edges: append([]string{}, task.Dependencies()...)})

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

func validateFinallyTasks(allTasksMap map[string]Task, finallyTaskStates map[string]state, finallyTaskNames []string) error {
	for _, taskName := range finallyTaskNames {
		if finallyTaskStates[taskName] == unvalidated {
			finallyTaskStates[taskName] = validating
			task, ok := allTasksMap[strings.ToLower(taskName)]
			if !ok {
				return fmt.Errorf("unknown task '%s'", taskName)
			}
			if len(task.Finally()) > 0 || hasNonOptionalArg(task) {
				return fmt.Errorf("'%s' not allowed in Finally", taskName)
			}
			if err := validateFinallyTasks(allTasksMap, finallyTaskStates, task.Dependencies()); err != nil {
				return err
			}
			finallyTaskStates[taskName] = valid
		} else if finallyTaskStates[taskName] == validating {
			return fmt.Errorf("Finally cycle detected")
		}
	}
	return nil
}

func hasNonOptionalArg(task Task) bool {
	for _, arg := range task.DeclaredArgs() {
		if arg.Validator != nil {
			return true
		}
	}
	return false
}
