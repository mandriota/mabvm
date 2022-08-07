//	Copyright 2022 Mark Mandriota
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package mabvm

import (
	"bytes"
	"runtime"
	"strconv"
	"sync/atomic"
)

type Code = byte

type Word = int64

type Status = Word

const (
	DIED Status = iota
	WORK
	RMUT
	WMUT
)

type Machine struct {
	Stat Status

	codP Word
	srcP Word
	dstP Word

	code []Code
	data []Word

	vtab []*Status

	// TODO: add debug mode
}

func (m *Machine) Init(code []Code, data []Word) {
	m.srcP = Word(len(data)) - 1
	m.dstP = Word(len(data)) - 1
	m.code = code
	m.data = data
}

func (m *Machine) Dump(dst []byte) []byte {
	b := bytes.NewBuffer(dst)

	b.WriteString("Stat=")
	b.WriteString(strconv.FormatInt(m.Stat, 16))
	b.WriteString("\ncodP=")
	b.WriteString(strconv.FormatInt(m.codP, 16))
	b.WriteString("\nsrcP=")
	b.WriteString(strconv.FormatInt(m.srcP, 16))
	b.WriteString("\ndstP=")
	b.WriteString(strconv.FormatInt(m.dstP, 16))

	for i, el := range m.data {
		b.WriteString("\ndata[")
		b.WriteString(strconv.FormatInt(int64(i), 16))
		b.WriteString("]=")
		b.WriteString(strconv.FormatInt(el, 16))
	}

	b.WriteByte('\n')

	return b.Bytes()
}

func (m *Machine) Run() {
	atomic.StoreInt64(&m.Stat, WORK)

	for m.codP = 0; m.codP < Word(len(m.code)); m.codP++ {
		op := m.code[m.codP]
		cc := Word(1)

		if op&EF == EF {
			cc = atomic.LoadInt64(&m.data[m.srcP])
			m.srcP--
		}

		srcD := atomic.LoadInt64(&m.data[m.srcP])
		dstD := atomic.LoadInt64(&m.data[m.dstP])

		if int64(op)&GC>>7*srcD >= dstD && int64(op)&EC>>6*srcD != dstD && int64(op)&LC>>5*srcD <= dstD {
			continue
		}

		if op&MF == MF {
			atomic.StoreInt64(&m.Stat, RMUT)

			for srcD == atomic.LoadInt64(&m.data[m.srcP]) {
				runtime.Gosched()
			}

			atomic.StoreInt64(&m.Stat, WORK)
		}

		// change cc sign if IF is setted
		cc -= int64(op) & IF >> 1 * cc

		switch op & JMask {
		case SJ:
			m.srcP += cc
		case DJ:
			m.dstP += cc
		case CJ:
			m.codP += cc
		case VJ:
			if p := m.vtab[m.dstP>>12]; p != &m.Stat {
				atomic.StoreInt64(&m.Stat, WMUT)

				for atomic.LoadInt64(p) != RMUT {
					runtime.Gosched()
				}

				atomic.StoreInt64(&m.Stat, WORK)
			}

			atomic.StoreInt64(&m.data[m.dstP], srcD+cc)
			m.srcP--
			m.dstP++
		}
	}

	atomic.StoreInt64(&m.Stat, DIED)
}
