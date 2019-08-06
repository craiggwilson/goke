package sh_test

import (
	"bytes"
	"context"

	"github.com/craiggwilson/goke/task"
)

func makeTestContext() *task.Context {
	var w bytes.Buffer
	return task.NewContext(context.Background(), &w, nil)
}