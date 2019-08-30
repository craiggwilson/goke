package task

// Task represents a task to be executed
type Task interface {
	DeclaredArgs() []DeclaredTaskArg
	Dependencies() []string
	Description() string
	Executor() Executor
	Hidden() bool
	Name() string
}

// DeclaredTaskArg is an argument for a particular task.
type DeclaredTaskArg struct {
	Name      string
	Required  bool
	Validator func(string) error
}
