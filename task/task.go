package task

// Task represents a task to be executed
type Task interface {
	Dependencies() []string
	Description() string
	Execute(*Context) error
	Name() string
}
