// SPDX-License-Identifier: GPL-3.0-or-later

package safewriter

import (
	"io"
	"sync"
)

func New(w io.Writer) io.Writer {
	return &Writer{
		mx: &sync.Mutex{},
		w:  w,
	}
}

type Writer struct {
	mx *sync.Mutex
	w  io.Writer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.mx.Lock()
	n, err = w.w.Write(p)
	w.mx.Unlock()
	return n, err
}
