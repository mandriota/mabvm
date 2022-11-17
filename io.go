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
)

// Writer - basic buffered writer interface.
// If last memory byte is 1 then flushes
// all memory without last bit to output.
type Writer struct {
	sync.RWMutex

	w io.Writer

	wmem []Word
	rmem []byte

	mtab *MutexTab
}

func NewWriter(w io.Writer, wm, rm []Word, mtab *MutexTab) *Writer {
	return &Writer{
		w:    w,
		wmem: wm,
		rmem: byteSliceOf(rm),
		mtab: mtab,
	}
}

func (w *Writer) Blocks() int {
	return (len(w.wmem) + 1) / BlockSize
}

func (w *Writer) Show() error {
	for {
		await(&w.RWMutex)

		if atomic.CompareAndSwapInt64(&w.wmem[len(w.wmem)-1], 1, 0) {
			wr := bufio.NewWriter(w.w)

			for i := 0; i < len(w.wmem)-1 && w.wmem[i+1] != 0; i += 2 {
				k := w.wmem[i+1]

				for j := w.wmem[i] / BlockSize; (j-1)*BlockSize <= k; j++ {
					(*w.mtab)[j].RLock()
					wr.Write(w.rmem[max(j*BlockSize, w.wmem[i])*8:][:min(k, (j+1)*BlockSize)*8])
					(*w.mtab)[j].RUnlock()
				}
			}

			wr.Flush()

			for i := range w.wmem {
				w.wmem[i] = 0
			}
		}
	}
}
