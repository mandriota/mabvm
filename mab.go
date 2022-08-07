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
	"strconv"
	"sync"
	"sync/atomic"
)

const VMBSize = 1 << 12

type Code = byte

type Word = int64

type Machine struct {
	sync.Mutex

	codP Word
	srcP Word
	dstP Word

	code []Code
	data []Word

	vtab []*sync.Mutex

	// TODO: add debug mode
}

func (mac *Machine) Init(code []Code, data []Word) {
	mac.srcP = Word(len(data)) - 1
	mac.dstP = Word(len(data)) - 1
	mac.code = code
	mac.data = data
	mac.vtab = make([]*sync.Mutex, len(data)/VMBSize+1)
}

func (mac *Machine) Dump(dst []byte) []byte {
	b := bytes.NewBuffer(dst)

	b.WriteString("codP=")
	b.WriteString(strconv.FormatInt(mac.codP, 16))
	b.WriteString("\nsrcP=")
	b.WriteString(strconv.FormatInt(mac.srcP, 16))
	b.WriteString("\ndstP=")
	b.WriteString(strconv.FormatInt(mac.dstP, 16))

	for i, el := range mac.data {
		b.WriteString("\ndata[")
		b.WriteString(strconv.FormatInt(int64(i), 16))
		b.WriteString("]=")
		b.WriteString(strconv.FormatInt(el, 16))
	}

	b.WriteByte('\n')

	return b.Bytes()
}

func (mac *Machine) Run() {
	for mac.codP = 0; mac.codP < Word(len(mac.code)); mac.codP++ {
		op := mac.code[mac.codP]
		cc := Word(1)

		if op&EF == EF {
			cc = atomic.LoadInt64(&mac.data[mac.srcP])
			mac.srcP--
		}

		srcD := atomic.LoadInt64(&mac.data[mac.srcP])
		dstD := atomic.LoadInt64(&mac.data[mac.dstP])

		if int64(op)&GC>>7*srcD >= dstD && int64(op)&EC>>6*srcD != dstD && int64(op)&LC>>5*srcD <= dstD {
			continue
		}

		if op&MF == MF && mac.TryLock() {
			mac.Lock()
		}

		// change cc sign if IF is setted
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

			if m := mac.vtab[mac.dstP/VMBSize]; m != nil && m != &mac.Mutex && !m.TryLock() {
				m.Unlock()
			}
		}
	}
}
