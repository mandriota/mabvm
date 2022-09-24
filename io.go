// Copyright 2022 Mark Mandriota
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// TODO: add flusher interface instead of writer

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
