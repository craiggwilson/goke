package task

// build begins building a task.
func build(name string) *Builder {
	task := &declaredTask{
		name: name,
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
func (b *Builder) Do(executor Executor) {
	b.task.executor = executor
}

// Hide the task from the task list.
func (b *Builder) Hide() {
	b.task.hidden = true
}

type declaredTask struct {
	name         string
	description  string
	dependencies []string
	executor     Executor
	hidden       bool
}

func (t *declaredTask) Dependencies() []string {
	return t.dependencies
}
func (t *declaredTask) Description() string {
	return t.description
}
func (t *declaredTask) Hidden() bool {
	return t.hidden
}
func (t *declaredTask) Executor() Executor {
	return t.executor
}
func (t *declaredTask) Name() string {
	return t.name
}
