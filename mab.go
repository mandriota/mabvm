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
		cx := Word(1)

		if op&EF == EF {
			cx = atomic.LoadInt64(&m.data[m.srcP])
			m.srcP--
		}

		srcD := atomic.LoadInt64(&m.data[m.srcP])
		dstD := atomic.LoadInt64(&m.data[m.dstP])

		// int64(op>>7)*srcD < dstD || int64(op>>6)*srcD == dstD || int64(op>>5)*srcD > dstD
		if !(((op&EC == 0) || (srcD == dstD)) && ((op&GC == 0) || (srcD > dstD)) && ((op&LC == 0) || (srcD < dstD))) {
			continue
		}

		if op&MF == MF {
			for srcD == atomic.LoadInt64(&m.data[m.srcP]) {
				runtime.Gosched()
			}
		}

		if op&IF == IF {
			cx *= -1
		}

		switch op & JMask {
		case SJ:
			m.srcP += cx
		case DJ:
			m.dstP += cx
		case CJ:
			m.codP += cx
		case VJ:
			atomic.StoreInt64(&m.data[m.dstP], srcD+cx)
			m.srcP--
			m.dstP++
		}
	}
}
