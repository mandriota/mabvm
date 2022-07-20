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
	"log"
	"runtime"
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

	l *log.Logger

	// TODO: add machine status public field
}

func NewMachine(l *log.Logger, code []Opcode, data []Word) *Machine {
	return &Machine{
		code: code,
		data: data,
		l:    l,
		srcP: int64(len(data)) - 1,
		dstP: int64(len(data)) - 1,
	}
}

func (m *Machine) dump() {
	m.l.Println("writing machine state dump ...")

	m.l.Printf("dump: codP=%d\n", m.codP)
	m.l.Printf("dump: srcP=%d\n", m.srcP)
	m.l.Printf("dump: dstP=%d\n", m.dstP)

	for i, el := range m.data {
		m.l.Printf("dump: data[%d]=%d\n", i, el)
	}
}

func (m *Machine) Run() {
	for m.codP = 0; m.codP < Word(len(m.code)); m.codP++ {
		op := m.code[m.codP]
		cx := Word(1)

		m.l.Printf("executing instruction %b (jump %d) ...\n", op, op&JMask)

		if op&EF == EF {
			cx = atomic.LoadInt64(&m.data[m.srcP])
			m.srcP--
		}

		srcD := atomic.LoadInt64(&m.data[m.srcP])
		dstD := atomic.LoadInt64(&m.data[m.dstP])

		if !(((op&EC == 0) || (srcD == dstD)) && ((op&GC == 0) || (srcD > dstD)) && ((op&LC == 0) || (srcD < dstD))) {
			m.l.Println("instruction not will executed: flags does not equal.")
			continue
		}

		if op&MF == MF {
			m.l.Println("interrupting machine: waiting for response ...")

			for srcD == atomic.LoadInt64(&m.data[m.srcP]) {
				runtime.Gosched()
			}
		}

		if op&IF == IF {
			cx *= -1
		}

		m.dump()

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

		m.dump()
	}
}
