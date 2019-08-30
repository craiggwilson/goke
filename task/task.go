package task

import "fmt"

// Task represents a task to be executed
type Task interface {
	DeclaredArgs() []DeclaredTaskArg
	Dependencies() []string
	Description() string
	Executor() Executor
	Hidden() bool
	Name() string
}

// Validator validates arguments.
type Validator func(string, string) error

// Required is a validator that ensures that an argument is present.
var Required = Validator(func(name, s string) error {
	if s == "" {
		return fmt.Errorf("argument %q is required, but was not supplied", name)
	}

	return nil
})

// ChainValidator is a validator that is the conjunction of the given validators.
func ChainValidator(validators ...Validator) Validator {
	return Validator(func(name, s string) error {
		for _, validator := range validators {
			if validator != nil {
				if err := validator(name, s); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// DeclaredTaskArg is an argument for a particular task.
type DeclaredTaskArg struct {
	Name      string
	Validator Validator
}
