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

func (b *Builder) Arg(a DeclaredTaskArg) *Builder {
	b.task.declaredArgs = append(b.task.declaredArgs, a)
	return b
}

// OptionalArg declares an option argument to the task.
func (b *Builder) OptionalArg(name string) *Builder {
	b.task.declaredArgs = append(b.task.declaredArgs, DeclaredTaskArg{
		Name:     name,
		Required: false,
	})
	return b
}

// RequiredArg declares an option argument to the task.
func (b *Builder) RequiredArg(name string) *Builder {
	b.task.declaredArgs = append(b.task.declaredArgs, DeclaredTaskArg{
		Name:     name,
		Required: true,
	})
	return b
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
func (b *Builder) Hide() *Builder {
	b.task.hidden = true
	return b
}

type declaredTask struct {
	name         string
	declaredArgs []DeclaredTaskArg
	description  string
	dependencies []string
	executor     Executor
	hidden       bool
}

func (t *declaredTask) DeclaredArgs() []DeclaredTaskArg {
	return t.declaredArgs
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
