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
	"bufio"
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
	rmem []byte
}

func NewWriter(w io.Writer, wm, rm []Word) *Writer {
	return &Writer{
		w:    w,
		wmem: wm,
		rmem: unsafe.Slice((*byte)(unsafe.Pointer(&rm[0])),
			uintptr(len(rm))*unsafe.Sizeof(rm[0])),
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
			wr := bufio.NewWriter(w.w)

			for i := 0; i < len(w.wmem)-1 && w.wmem[i+1] != 0; i += 2 {
				wr.Write(w.rmem[w.wmem[i]*8 : (w.wmem[i]+w.wmem[i+1])*8])
			}

			wr.Flush()

			for i := range w.wmem {
				w.wmem[i] = 0
			}
		}
	}
}
