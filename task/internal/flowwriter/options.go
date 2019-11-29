package flowwriter

import "strings"

// Option configures a Writer.
type Option func(*Writer)

// WithCutset sets the runes responsible for allowing a line to break.
func WithCutset(cutset string) Option {
	return WithWrapFunc(func(r rune) bool {
		return strings.ContainsRune(cutset, r)
	})
}

// WithWrapFunc sets the func responsible for allowing a line to break.
func WithWrapFunc(f func(rune) bool) Option {
	return func(w *Writer) {
		w.canWrap = f
	}
}

// WithIndent sets the text to prepend to a wrapped line.
func WithIndent(indent []byte) Option {
	return func(w *Writer) {
		w.indent = indent
	}
}

// WrapAtColumn indicates to a writer to wrap text at the column.
func WrapAtColumn(column int) Option {
	return func(w *Writer) {
		w.maxLineLen = column
	}
}
