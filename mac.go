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
	"io"
	"strconv"
	"sync"
	"sync/atomic"
)

type Binder interface {
	Bind(m *sync.Mutex, blocks int)
}

const BlockSize = 1 << 12

type Code = byte

type Word = int64

type MutexTab []*sync.Mutex

type Machine struct {
	sync.Mutex

	codP Word
	srcP Word
	dstP Word

	code []Code
	data []Word

	mtab *MutexTab
}

func NewMachine(code []Code, data []Word, mtab *MutexTab) *Machine {
	mac := &Machine{
		srcP: Word(len(data)) - 1,
		code: code,
		data: data,
		mtab: mtab,
	}

	mac.Bind(&mac.Mutex, len(mac.data)>>12-len(*mac.mtab)+1)

	return mac
}

func (mac *Machine) Bind(m *sync.Mutex, blocks int) {
	*mac.mtab = append(*mac.mtab, make([]*sync.Mutex, blocks)...)

	for i := range (*mac.mtab)[len(*mac.mtab)-blocks:] {
		(*mac.mtab)[i] = m
	}
}

func (mac *Machine) dumpFlag(w *bufio.Writer, f byte, msg string) {
	if mac.code[mac.codP-1]&f == f {
		w.WriteString(msg)
	}
}

func (mac *Machine) Dump(w *bufio.Writer) {
	w.WriteString("\n==============================")

	switch mac.code[mac.codP-1] & JMask {
	case SJ:
		w.WriteString("\nSJ: Source Jump")
	case DJ:
		w.WriteString("\nDJ: Destination Jump")
	case CJ:
		w.WriteString("\nCJ: Code Jump")
	case VJ:
		w.WriteString("\nVJ: Value Jump")
	}

	w.WriteString("\nFlags:")

	mac.dumpFlag(w, IF, "\n\tIF: Inversion Flag")
	mac.dumpFlag(w, EF, "\n\tEF: Extension Flag")
	mac.dumpFlag(w, MF, "\n\tMF: Mutex Flag")
	mac.dumpFlag(w, LC, "\n\tLC: Lower Conditional")
	mac.dumpFlag(w, EC, "\n\tEC: Equal Conditional")
	mac.dumpFlag(w, GC, "\n\tGC: Greater Conditional")

	w.WriteString("\ncodP: Code Pointer: ")
	w.WriteString(strconv.FormatInt(mac.codP, 16))
	w.WriteString("\nsrcP: Source Pointer: ")
	w.WriteString(strconv.FormatInt(mac.srcP, 16))
	w.WriteString("\ndstP: Destination Pointer: ")
	w.WriteString(strconv.FormatInt(mac.dstP, 16))

	w.WriteString("\nData:")

	for i, el := range mac.data {
		w.WriteString("\n\tWord[")
		w.WriteString(strconv.FormatInt(int64(i), 16))
		w.WriteString("]: ")
		w.WriteString(strconv.FormatInt(el, 16))
	}

	w.WriteString("\n==============================\n")
}

func (mac *Machine) Tick() {
	op := mac.code[mac.codP]

	if op&MF == MF {
		await(&mac.Mutex)
	}

	cc := Word(1)

	if op&EF == EF {
		cc = atomic.LoadInt64(&mac.data[mac.srcP])
		mac.srcP--
	}

	srcD := atomic.LoadInt64(&mac.data[mac.srcP])
	dstD := atomic.LoadInt64(&mac.data[mac.dstP])

	if int64(op)&GC>>7*srcD >= dstD &&
		int64(op)&EC>>6*srcD != dstD &&
		int64(op)&LC>>5*srcD <= dstD {
		return
	}

	cc -= int64(op) & IF >> 1 * cc

	switch op & JMask {
	case SJ:
		mac.srcP += cc
	case DJ:
		mac.dstP += cc
	case CJ:
		mac.codP += cc
	case VJ:
		atomic.StoreInt64(&mac.data[mac.dstP], srcD+cc)
		mac.srcP--
		mac.dstP++

		m := (*mac.mtab)[mac.dstP/BlockSize]
		if m != nil && m != &mac.Mutex {
			m.TryLock()
		}
	}

	mac.codP++
}

func (mac *Machine) Show() {
	mac.codP = 0

	for mac.codP < Word(len(mac.code)) {
		mac.Tick()
	}
}

func (mac *Machine) DebugShow(dw io.Writer) {
	w := bufio.NewWriter(dw)
	defer w.Flush()

	mac.codP = 1
	mac.Dump(w)
	mac.codP = 0

	for mac.codP < Word(len(mac.code)) {
		mac.Tick()
		mac.Dump(w)
	}
}
