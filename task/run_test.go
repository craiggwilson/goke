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
	finally          []string
	shouldError      bool
	shouldFailRun    bool
	expectedRunOrder []string
}

func runTests(t *testing.T, registry *Registry, declareTasks bool, testCases []testTaskCfg) {
	for _, tc := range testCases {
		runOrder = []string{}
		if declareTasks {
			declare(registry, tc.name, tc.shouldError).Finally(tc.finally...)
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

func TestFinally(t *testing.T) {
	dummy := []string{}
	reg := NewRegistry()

	// Finally() allows simple tasks and aggregates thereof
	declare(reg, "ok1", false)
	declare(reg, "ok2", false)
	declare(reg, "err", true)
	declare(reg, "dep", false).DependsOn("ok1", "ok2")

	// Finally() does *not* allow tasks which
	// - have required arguments
	declare(reg, "requiredArg", false).RequiredArg("foo")
	// - use Finally() themselves
	declare(reg, "hasFinally", false).Finally("ok1")
	// - have dependencies that are not allowed
	declare(reg, "dependsOnNotAllowed", false).DependsOn("requiredArg")

	t.Run("NotAllowedInPureAggregates", func(t *testing.T) {
		reg.Declare("t0").DependsOn("ok1", "ok2").Finally("ok1")
		runTests(t, reg, false, []testTaskCfg{
			{"t0", dummy, false, true, dummy},
		})
	})

	t.Run("ShouldRunOnTaskSuccess", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t1", []string{"ok1"}, false, false, []string{"t1", "ok1"}},
			{"t2", []string{"ok1", "ok1", "ok2"}, false, false, []string{"t2", "ok1", "ok2"}},
			{"t3", []string{"dep"}, false, false, []string{"t3", "ok1", "ok2", "dep"}},
		})
	})

	t.Run("ShouldRunOnTaskFailure", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t4", []string{"ok1"}, true, true, []string{"t4", "ok1"}},
			{"t5", []string{"ok1", "ok2"}, true, true, []string{"t5", "ok1", "ok2"}},
		})
	})

	t.Run("ShouldIgnoreErrorsInFinally", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t6", []string{"err", "ok1"}, false, false, []string{"t6", "err", "ok1"}},
		})
	})

	t.Run("ShouldErrorIfMisconfigured", func(t *testing.T) {
		runTests(t, reg, true, []testTaskCfg{
			{"t7", []string{"requiredArg"}, false, true, dummy},
			{"t8", []string{"hasFinally"}, false, true, dummy},
			{"t9", []string{"dependsOnNotAllowed"}, false, true, dummy},
			{"t10", []string{"doesNotExist"}, false, true, dummy},
		})
	})

	t.Run("ShouldReverseFinalizeDependencies", func(t *testing.T) {
		declare(reg, "t11", false).DependsOn("t1", "t2").Finally("err")
		declare(reg, "t12", true).DependsOn("t1", "t2").Finally("dep")
		declare(reg, "t13", false).DependsOn("t11").Finally("ok1")
		declare(reg, "t14", false).DependsOn("t12").Finally("ok1")

		runTests(t, reg, false, []testTaskCfg{
			{"t11", dummy, false, false, []string{"t1", "t2", "t11", "err", "ok1", "ok2"}},
			{"t12", dummy, false, true, []string{"t1", "t2", "t12", "ok1", "ok2", "dep"}},
			{"t13", dummy, false, false, []string{"t1", "t2", "t11", "t13", "ok1", "err", "ok2"}},
			{"t14", dummy, false, true, []string{"t1", "t2", "t12", "ok1", "ok2", "dep"}},
		})
	})
}
