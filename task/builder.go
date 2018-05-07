package task

import "fmt"

// build begins building a task.
func build(name string) *Builder {
	task := &declaredTask{
		name: name,
		execute: func(*Context) error {
			return fmt.Errorf("no executor defined")
		},
	}
	return &Builder{task: task}
}

// Builder provides a fluent way to build up a task.
type Builder struct {
	task *declaredTask
}

// Description sets the description for the task.
func (b *Builder) Description(description string) *Builder {
	b.task.description = description
	return b
}

// DependsOn declares other tasks which must run before this one.
func (b *Builder) DependsOn(names ...string) *Builder {
	b.task.dependencies = names
	return b
}

// Do declares the executor when this task runs.
func (b *Builder) Do(execute func(*Context) error) {
	b.task.execute = execute
}

type declaredTask struct {
	name         string
	description  string
	dependencies []string
	execute      func(*Context) error
}

func (t *declaredTask) Dependencies() []string {
	return t.dependencies
}
func (t *declaredTask) Description() string {
	return t.description
}
func (t *declaredTask) Execute(ctx *Context) error {
	return t.execute(ctx)
}
func (t *declaredTask) Name() string {
	return t.name
}
