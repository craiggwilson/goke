package flowwriter

import (
	"bytes"
	"io"
	"unicode"
)

// New makes a Writer.
func New(w io.Writer, opts ...Option) *Writer {
	wr := &Writer{
		w:          w,
		maxLineLen: 120,
		indent:     nil,
		canWrap:    unicode.IsSpace,
	}

	for _, opt := range opts {
		opt(wr)
	}

	return wr
}

// Writer wraps text at a specified column width.
type Writer struct {
	w          io.Writer
	maxLineLen int
	indent     []byte
	canWrap    func(rune) bool

	buf     bytes.Buffer
	wordBuf bytes.Buffer
	lineLen int
}

// Write implements the io.Writer interface.
func (w *Writer) Write(p []byte) (n int, err error) {
	var nn int
	for _, r := range string(p) {
		nn, err = w.wordBuf.WriteRune(r)
		n += nn
		if err != nil {
			return
		}

		if r != '\n' && !w.canWrap(r) {
			continue
		}

		if w.lineLen+w.wordBuf.Len() > w.maxLineLen {
			if err = w.buf.WriteByte('\n'); err != nil {
				return
			}

			if len(w.indent) > 0 {
				if _, err = w.buf.Write(w.indent); err != nil {
					return
				}
			}
			w.lineLen = len(w.indent)
		}

		w.lineLen += w.wordBuf.Len()
		if _, err = w.wordBuf.WriteTo(&w.buf); err != nil {
			return
		}

		if r == '\n' {
			w.lineLen = 0
		}
	}

	if w.buf.Len() > 0 {
		w.buf.WriteTo(w.w)
	}

	return n, nil
}
