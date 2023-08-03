package task

import (
	"errors"
	"testing"
)

func TestTaskWithFinally(t *testing.T) {
	finalizedCount := 0

	finalizer := func(ctx *Context) error {
		finalizedCount++
		return nil
	}
	finalizerWithError := func(ctx *Context) error {
		return errors.New("uh oh")
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
		// Error in finalizer should not affect task success
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
			t.Error("expected an error")
		} else if !tc.taskError && err != nil {
			t.Error("expected no error")
		}
		if finalizedCount != tc.expectedFinalizedCount {
			t.Errorf("expected %d finalizer(s) to run but got %d", tc.expectedFinalizedCount, finalizedCount)
		}
	}
}

func TestShouldExitOnAggregateTaskWithFinalizer(t *testing.T) {
	shouldNeverRun := func(ctx *Context) error {
		return errors.New("task or finalizer ran")
	}
	registry := NewRegistry()
	registry.Declare("actual_task").Do(shouldNeverRun)
	registry.Declare("aggregate_task").DependsOn("actual_task").Finally(shouldNeverRun)

	if err := Run(registry, []string{"aggregate_task"}); err == nil {
		t.Error("expected an error")
	} else if err.Error() == "task(s) [actual_task] failed" {
		t.Error("task should not have been executed")
	}
}
