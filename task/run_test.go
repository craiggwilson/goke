package task

import (
	"fmt"
	"reflect"
	"testing"
)

var runOrder []string

func makeExecutor(name string, shouldError bool) Executor {
	return func(ctx *Context) error {
		runOrder = append(runOrder, name)
		if shouldError {
			return fmt.Errorf("error in %s", name)
		}
		return nil
	}
}

func declare(registry *Registry, name string, shouldError bool) *Builder {
	b := registry.Declare(name)
	b.Do(makeExecutor(name, shouldError))
	return b
}

type testTaskCfg struct {
	name             string
	deferredTasks    []string
	shouldError      bool
	shouldFailRun    bool
	expectedRunOrder []string
}

func runTests(t *testing.T, registry *Registry, declareTasks bool, testCases []testTaskCfg) {
	for _, tc := range testCases {
		runOrder = []string{}
		if declareTasks {
			declare(registry, tc.name, tc.shouldError).Defer(tc.deferredTasks...)
		}
		err := Run(registry, []string{tc.name})
		if err == nil && tc.shouldFailRun {
			t.Fatalf("expected error")
		} else if err != nil && !tc.shouldFailRun {
			t.Fatalf("expected no error")
		}
		if !reflect.DeepEqual(runOrder, tc.expectedRunOrder) {
			t.Fatalf("case '%s': expected run order %v but got %v", tc.name, tc.expectedRunOrder, runOrder)
		}
	}
}

func TestDefer(t *testing.T) {
	dummy := []string{}
	reg := NewRegistry()

	// Defer() allows most *almost* any task...
	declare(reg, "ok1", false)
	declare(reg, "ok2", false)
	declare(reg, "err", true)
	declare(reg, "dep", false).DependsOn("ok1", "ok2")
	declare(reg, "requiredArg", false).RequiredArg("foo")
	reg.Declare("agg").DependsOn("ok1", "ok2")

	// ...unless they use Defer() themselves or depend on such a task
	declare(reg, "hasDefer", false).Defer("ok1")
	declare(reg, "dependsOnHasDefer", false).DependsOn("hasDefer")

	t.Run("ShouldRunOnTaskSuccess", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t1", []string{"ok1"}, false, false, []string{"t1", "ok1"}},
			{"t2", []string{"ok1", "ok1", "ok2"}, false, false, []string{"t2", "ok1", "ok2"}},
			{"t3", []string{"dep"}, false, false, []string{"t3", "ok1", "ok2", "dep"}},
		})
	})

	t.Run("ShouldRunOnTaskFailure", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t4", []string{"agg"}, true, true, []string{"t4", "ok1", "ok2"}},
			{"t5", []string{"ok1", "ok2"}, true, true, []string{"t5", "ok1", "ok2"}},
		})
	})

	t.Run("ShouldIgnoreErrorsInDefer", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t6", []string{"err", "ok1"}, false, false, []string{"t6", "err", "ok1"}},
			// if a deferred task is missing a required argument, it should *not* execute
			{"t7", []string{"requiredArg", "ok1"}, false, false, []string{"t7", "ok1"}},
		})
	})

	t.Run("ShouldErrorIfMisconfigured", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t8", []string{"hasDefer"}, false, true, dummy},
			{"t9", []string{"dependsOnHasDefer"}, false, true, dummy},
			{"t10", []string{"doesNotExist"}, false, true, dummy},
		})
	})

	t.Run("ShouldReverseFinalizeDependencies", func(t *testing.T) {
		declare(reg, "t11", false).DependsOn("t1", "t2").Defer("err")
		declare(reg, "t12", true).DependsOn("t1", "t2").Defer("dep")
		declare(reg, "t13", false).DependsOn("t11").Defer("ok1")
		declare(reg, "t14", false).DependsOn("t12").Defer("ok1")

		runTests(t, reg, false, []testTaskCfg{
			{"t11", dummy, false, false, []string{"t1", "t2", "t11", "err", "ok1", "ok2"}},
			{"t12", dummy, false, true, []string{"t1", "t2", "t12", "ok1", "ok2", "dep"}},
			{"t13", dummy, false, false, []string{"t1", "t2", "t11", "t13", "ok1", "err", "ok2"}},
			{"t14", dummy, false, true, []string{"t1", "t2", "t12", "ok1", "ok2", "dep"}},
		})
	})
}
