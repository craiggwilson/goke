package task

// Task represents a task to be executed
type Task interface {
	Dependencies() []string
	Description() string
	Executor() Executor
	Hidden() bool
	Name() string
}
