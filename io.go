package mabvm

import (
	"io"
	"sync"
	"unsafe"
)

type Writer struct {
	sync.Mutex

	w io.Writer

	rdata []byte
	ldata []byte
}

func NewWriter(w io.Writer, d []Word) *Writer {
	writer := &Writer{
		w: w,
		rdata: unsafe.Slice(
			(*byte)(unsafe.Pointer(&d[0])),
			uintptr(len(d))*unsafe.Sizeof(d[0]),
		),
	}
	writer.ldata = make([]byte, len(writer.rdata))
	return writer
}

func (w *Writer) Blocks() int {
	return len(w.rdata)*int(unsafe.Sizeof(Word(0)))/BlockSize + 1
}

func (w *Writer) Run() error {
	for {
		if w.TryLock() {
			copy(w.ldata, w.rdata)
			w.Lock()
		} else {
			w.Unlock()
		}

		j := 0

		for i, el := range w.rdata {
			if el != w.ldata[i] {
				w.ldata[j] = el
				j++
			}
		}

		w.w.Write(w.ldata[:j])
	}
}
