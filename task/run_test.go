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
	finalizerWithError := func(ctx *Context) error {
		return errors.New("uh oh...")
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
		// Error in finalizer
		{[]Executor{finalizerWithError}, false, true, 0},
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

		err := Run(registry, []string{"foo"})
		if tc.taskError && err == nil {
			t.Error("Expected an error")
		} else if !tc.taskError && err != nil {
			t.Error("Expected no error")
		}
		if finalizedCount != tc.expectedFinalizedCount {
			t.Errorf("Expected %d finalizer(s) to run but got %d", tc.expectedFinalizedCount, finalizedCount)
		}
	}
}

func TestShouldExitOnAggregateTaskWithFinalizer(t *testing.T) {
	doNothing := func(ctx *Context) error {
		return nil
	}
	registry := NewRegistry()
	registry.Declare("actual_task").Do(doNothing)
	registry.Declare("aggregate_task").DependsOn("first_task").Finally(doNothing)

	err := Run(registry, []string{"aggregate_task"})
	if err == nil {
		t.Error("Expected error")
	}
}
