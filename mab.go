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

// TODO: add multithreading support

type Opcode = byte

type Word = int64 // machine word

type Machine struct {
	code []Opcode
	data []Word

	codP Word
	srcP Word
	dstP Word

	// TODO: add machine status public field
	// TODO: add debug mode
}

func NewMachine(code []Opcode, data []Word) *Machine {
	return &Machine{
		code: code,
		data: data,
		srcP: int64(len(data)) - 1,
		dstP: int64(len(data)) - 1,
	}
}

func (m *Machine) Dump(dst []byte) []byte {
	b := bytes.NewBuffer(dst)

	b.WriteString("codP=")
	b.WriteString(strconv.FormatInt(m.codP, 10))
	b.WriteString("\nsrcP=")
	b.WriteString(strconv.FormatInt(m.srcP, 10))
	b.WriteString("\ndstP=")
	b.WriteString(strconv.FormatInt(m.dstP, 10))

	for i, el := range m.data {
		b.WriteString("\ndata[")
		b.WriteString(strconv.FormatInt(int64(i), 10))
		b.WriteString("]=")
		b.WriteString(strconv.FormatInt(el, 10))
	}

	b.WriteByte('\n')

	return b.Bytes()
}

func (m *Machine) Run() {
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
			for srcD == atomic.LoadInt64(&m.data[m.srcP]) {
				runtime.Gosched()
			}
		}

		// change cx sign if IF is setted
		cc -= int64(op) & IF >> 1 * cc

		switch op & JMask {
		case SJ:
			m.srcP += cc
		case DJ:
			m.dstP += cc
		case CJ:
			m.codP += cc
		case VJ:
			atomic.StoreInt64(&m.data[m.dstP], srcD+cc)
			m.srcP--
			m.dstP++
		}
	}
}
