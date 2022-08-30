package mabvm

import (
	_ "embed"
	"io"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Writer - basic buffered writer interface.
// If last memory byte is 1 then flushes
// all memory without last bit to output.
type Writer struct {
	sync.Mutex

	w io.Writer

	wmem []Word
	bmem []byte
}

func NewWriter(w io.Writer, d []Word) *Writer {
	return &Writer{
		w:    w,
		wmem: d,
		bmem: unsafe.Slice(
			(*byte)(unsafe.Pointer(&d[0])),
			uintptr(len(d))*unsafe.Sizeof(d[0]),
		),
	}
}

func (w *Writer) Blocks() int {
	return (len(w.wmem) + 1) / BlockSize
}

func (w *Writer) Show() error {
	for {
		if w.TryLock() {
			w.Lock()
		} else {
			w.Unlock()
		}

		if atomic.CompareAndSwapInt64(&w.wmem[len(w.wmem)-1], 1, 0) {
			if _, err := w.w.Write(w.bmem[:len(w.bmem)-8]); err != nil {
				return err
			}

			for i := range w.wmem {
				w.wmem[i] = 0
			}
		}
	}
}
