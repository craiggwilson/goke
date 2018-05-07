package task

// NewRegistry creates a new registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Registry holds all the tasks able to be run.
type Registry struct {
	tasks []Task
}

// Register a task in the Configuration.
func (r *Registry) Register(task Task) {
	r.tasks = append(r.tasks, task)
}

// Declare a task to be registered.
func (r *Registry) Declare(name string) *Builder {
	tb := build(name)
	r.Register(tb.task)
	return tb
}
