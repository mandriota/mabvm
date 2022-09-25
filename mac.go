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
	"bytes"
	"os"
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

type Machine struct {
	sync.Mutex

	codP Word
	srcP Word
	dstP Word

	code []Code
	data []Word

	mtab []*sync.Mutex
}

func NewMachine(code []Code, data []Word, mtab []*sync.Mutex) *Machine {
	mac := &Machine{
		srcP: Word(len(data)) - 1,
		code: code,
		data: data,
		mtab: mtab,
	}

	mac.Bind(&mac.Mutex, len(mac.data)>>12-len(mac.mtab)+1)

	return mac
}

func (mac *Machine) Bind(m *sync.Mutex, blocks int) {
	mac.mtab = append(mac.mtab, make([]*sync.Mutex, blocks)...)

	for i := range mac.mtab[len(mac.mtab)-blocks:] {
		mac.mtab[i] = m
	}
}

func (mac *Machine) Dump(dst []byte) []byte {
	b := bytes.NewBuffer(dst)
	b.WriteString("\n====================")

	c := mac.code[mac.codP]

	switch c & JMask {
	case SJ:
		b.WriteString("\nSource Jump (SJ)")
	case DJ:
		b.WriteString("\nDestination Jump (DJ)")
	case CJ:
		b.WriteString("\nCode Jump (CJ)")
	case VJ:
		b.WriteString("\nValue Jump (VJ)")
	}

	b.WriteString("\nFlags:")

	if c&IF == IF {
		b.WriteString("\n\tIF")
	}

	if c&EF == EF {
		b.WriteString("\n\tEF")
	}

	if c&MF == MF {
		b.WriteString("\n\tMF")
	}

	if c&LC == LC {
		b.WriteString("\n\tLC")
	}

	if c&EC == EC {
		b.WriteString("\n\tEC")
	}

	if c&GC == GC {
		b.WriteString("\n\tGC")
	}

	b.WriteString("\ncodP: ")
	b.WriteString(strconv.FormatInt(mac.codP, 16))
	b.WriteString("\nsrcP: ")
	b.WriteString(strconv.FormatInt(mac.srcP, 16))
	b.WriteString("\ndstP: ")
	b.WriteString(strconv.FormatInt(mac.dstP, 16))

	b.WriteString("\ndata:")

	for i, el := range mac.data {
		b.WriteString("\n\tword[")
		b.WriteString(strconv.FormatInt(int64(i), 16))
		b.WriteString("]: ")
		b.WriteString(strconv.FormatInt(el, 16))
	}

	b.WriteString("\n====================\n")

	return b.Bytes()
}

func (mac *Machine) Tick() {
	mac.codP++

	op := mac.code[mac.codP]

	if op&MF == MF {
		if mac.TryLock() {
			mac.Lock()
		} else {
			mac.Unlock()
		}
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

		m := mac.mtab[mac.dstP/BlockSize]
		if m != nil && m != &mac.Mutex && !m.TryLock() {
			m.Unlock()
		}
	}
}

func (mac *Machine) Show() {
	mac.codP = -1

	for mac.codP+1 < Word(len(mac.code)) {
		mac.Tick()
	}
}

func (mac *Machine) DebugShow() {
	buf := make([]byte, 0, 4096)

	w := bufio.NewWriter(os.Stdout)

	buf = mac.Dump(buf)
	w.Write(buf)
	buf = buf[:0]

	mac.codP = -1

	for mac.codP+1 < Word(len(mac.code)) {
		mac.Tick()
		buf = mac.Dump(buf)
		w.Write(buf)
		buf = buf[:0]
	}

	w.Flush()
}
