package task

import (
	"testing"
)

func TestUnusedArgs(t *testing.T) {
	registry := NewRegistry()

	fooArgs := []string{"bar", "baz"}
	registry.Declare("foo").Description("the foo command will do nothing").OptionalArgs(fooArgs...).Do(func(ctx *Context) error {
		return nil
	})

	quuxArgs := []string{"quuz", "corge"}
	registry.Declare("quux").Description("the quux command will do nothing").OptionalArgs(quuxArgs...).Do(func(ctx *Context) error {
		return nil
	})

	unusedArgs := getUnusedArgs(registry.Tasks(), globalArgs{
		"": {
			"bar":    "bar",
			"baz":    "baz",
			"quuz":   "quuz",
			"unused": "unused",
		},
		"quux": {
			"corge": "corge",
			"fake":  "fake",
		},
		"invalidTest": {
			"bar": "bar",
		},
	})

	// Using a map because there is no inherent order that unusedArgs should be.
	expectedUnusedArgs := map[string]bool{"unused": true, "quux:fake": true, "invalidTest:bar": true}

	if len(unusedArgs) != len(expectedUnusedArgs) {
		t.Fatalf("expected args length of %d, instead got %d", len(expectedUnusedArgs), len(unusedArgs))
	} else {
		for _, arg := range unusedArgs {
			if _, ok := expectedUnusedArgs[arg]; !ok {
				t.Fatalf("got unexpected arg %s", arg)
			}
		}
	}
}

func TestShouldExitOnUnusedArgs(t *testing.T) {
	registryErrorOnUnused := NewRegistry(WithShouldErrorOnUnusedArgs(true))
	registryNoErrorOnUnused := NewRegistry(WithShouldErrorOnUnusedArgs(false))

	fooArgs := []string{"bar", "baz"}
	registryErrorOnUnused.Declare("foo").Description("the foo command will do nothing").OptionalArgs(fooArgs...).Do(func(ctx *Context) error {
		return nil
	})
	registryNoErrorOnUnused.Declare("foo").Description("the foo command will do nothing").OptionalArgs(fooArgs...).Do(func(ctx *Context) error {
		return nil
	})

	quuxArgs := []string{"quuz", "corge"}
	registryErrorOnUnused.Declare("quux").Description("the quux command will do nothing").OptionalArgs(quuxArgs...).Do(func(ctx *Context) error {
		return nil
	})
	registryNoErrorOnUnused.Declare("quux").Description("the quux command will do nothing").OptionalArgs(quuxArgs...).Do(func(ctx *Context) error {
		return nil
	})

	testCases := []struct {
		args        []string
		shouldError bool
	}{
		{[]string{"foo", "--bar", "--baz", "quux"}, false},
		{[]string{"foo", "--bar", "--baz", "quux", "--fake"}, true},
		{[]string{"foo", "--bar", "--baz", "quux", "--quux:corge"}, false},
		{[]string{"foo", "--bar", "--baz", "quux", "--quux:fake"}, true},
	}

	for _, tc := range testCases {
		err := Run(registryErrorOnUnused, tc.args)
		if tc.shouldError && err == nil {
			t.Errorf("expecting an error")
		} else if !tc.shouldError && err != nil {
			t.Errorf("not expecting an error")
		}

		err = Run(registryNoErrorOnUnused, tc.args)
		if err != nil {
			t.Errorf("not expecting an error")
		}
	}
}
