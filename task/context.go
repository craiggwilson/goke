package task

import (
	"context"
	"fmt"
	"io"
)

// Context holds information relevent to executing tasks.
type Context struct {
	context.Context

	Args    []string
	DryRun  bool
	Verbose bool

	w io.Writer
}

// Log formats using the default formats for its operands sends it to the log.
// Spaces are added between operands when neither is a string.
func (ctx *Context) Log(v ...interface{}) {
	fmt.Fprint(ctx.w, v...)
}

// Logln formats using the default formats for its operands and sends it to the log.
// Spaces are always added between operands and a newline is appended.
func (ctx *Context) Logln(v ...interface{}) {
	fmt.Fprintln(ctx.w, v...)
}

// Logf formats according to a format specifier and sends it to the log.
func (ctx *Context) Logf(format string, v ...interface{}) {
	fmt.Fprintf(ctx.w, format, v...)
}

// Writer implements the io.Writer interface.
func (ctx *Context) Write(p []byte) (n int, err error) {
	return ctx.w.Write(p)
}
