package task

import (
	"fmt"
	"testing"
)

func TestToposort(t *testing.T) {

	t.Run("Should sort correctly", func(t *testing.T) {
		g := []*graphNode{
			{
				task:  dummyTask("5"),
				edges: []string{"11"},
			},
			{
				task:  dummyTask("7"),
				edges: []string{"11", "8"},
			},
			{
				task:  dummyTask("3"),
				edges: []string{"8", "10"},
			},
			{
				task:  dummyTask("11"),
				edges: []string{"2", "9", "10"},
			},
			{
				task:  dummyTask("8"),
				edges: []string{"9"},
			},
			{
				task: dummyTask("2"),
			},
			{
				task: dummyTask("9"),
			},
			{
				task: dummyTask("10"),
			},
		}

		result, err := toposort(g)
		if err != nil {
			t.Fatalf("expected no error, but got %s", err)
		}

		if fmt.Sprint(result) != "[2 9 10 8 11 3 5 7]" {
			t.Fatalf("expected [2 9 10 8 11 3 5 7], but got %s", result)
		}
	})
	t.Run("Should error on a cycle", func(t *testing.T) {
		g := []*graphNode{
			{
				task:  dummyTask("7"),
				edges: []string{"11"},
			},
			{
				task:  dummyTask("11"),
				edges: []string{"10"},
			},
			{
				task:  dummyTask("10"),
				edges: []string{"7"},
			},
		}

		_, err := toposort(g)
		if err == nil {
			t.Fatal("expected an error, but got none")
		}
	})

}

type dummyTask string

func (t dummyTask) Dependencies() []string {
	return nil
}
func (t dummyTask) Description() string {
	return ""
}
func (t dummyTask) Execute(*Context) error {
	return nil
}
func (t dummyTask) Name() string {
	return string(t)
}
