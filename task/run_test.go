package task

import (
	"errors"
	"testing"
)

func TestTaskWithFinalizer(t *testing.T) {
	finalizedCount := 0

	finalizer := func(ctx *Context) error {
		finalizedCount++
		return nil
	}

	testCases := []struct {
		finalizers             []Executor
		taskError              bool
		finalizeOnError        bool
		expectedFinalizedCount int
	}{
		// No finalizeer
		{[]Executor{}, false, true, 0},
		// Single finalizer
		{[]Executor{finalizer}, false, true, 1},
		// Multiple finalizers
		{[]Executor{finalizer, finalizer, finalizer}, false, true, 3},
		// Should finalize by default if task had error
		{[]Executor{finalizer}, true, true, 1},
		// SkipFinallyOnError should prevent finalizer execution if task had error
		{[]Executor{finalizer}, true, false, 0},
	}

	for _, tc := range testCases {
		finalizedCount = 0
		registry := NewRegistry()
		b := registry.Declare("foo").Finally(tc.finalizers...)
		if !tc.finalizeOnError {
			b.SkipFinallyOnError()
		}
		b.Do(func(ctx *Context) error {
			if tc.taskError {
				return errors.New("task failed")
			}
			return nil
		})

		Run(registry, []string{"foo"})
		if finalizedCount != tc.expectedFinalizedCount {
			t.Errorf("Expected %d finalizer(s) to run but got %d", tc.expectedFinalizedCount, finalizedCount)
		}
	}
}
