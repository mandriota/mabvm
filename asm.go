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
	"fmt"
	"io"
)

type AsmParser struct {
	src string
	pos int
}

//TODO: Add buildError function instead of fmt.Errorf

func (ap *AsmParser) Reset(src string) {
	ap.src = src
	ap.pos = 0
}

func (ap *AsmParser) readByte() {
	ap.pos++
}

func (ap *AsmParser) peekByte() byte {
	if ap.pos < len(ap.src) {
		return ap.src[ap.pos]
	}

	return 0
}

func (ap *AsmParser) Parse(mac *Machine) error {
	if mac == nil {
		mac = &Machine{}
		mac.Init(make([]byte, 0, 1<<16), make([]int64, 0, 1<<13))
	}

	for {
		switch err := ap.parseOpcode(mac); err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func (ap *AsmParser) parseNumber(mac *Machine) (err error) {
	sign := int64(0)
	word := int64(0)
	base := int64(0)

	switch ap.peekByte() {
	case '+':
		sign = +1
	case '-':
		sign = -1
	case '\x00':
		return io.EOF
	default:
		return fmt.Errorf("unexpected character")
	}

	ap.readByte()

	switch ap.peekByte() {
	case 'b':
		base = 2
	case 'o':
		base = 8
	case 'd':
		base = 10
	case '\x00':
		return fmt.Errorf("unexpected EOF")
	default:
		return fmt.Errorf("unexpected character")
	}

	ap.readByte()

	for c := ap.peekByte(); c >= '0' && c < '0'+byte(base); c = ap.peekByte() {
		word *= base
		word += int64(c - '0')
		ap.readByte()
	}

	mac.data = append(mac.data, word*sign)
	return nil
}

func (ap *AsmParser) parseOpcode(mac *Machine) (err error) {
	ap.nextSect()

	if ap.peekByte() != ':' {
		return ap.parseNumber(mac)
	}

	ap.readByte()

	op := Code(0)

	switch ap.peekByte() {
	case 'S':
		op = SJ
	case 'D':
		op = DJ
	case 'C':
		op = CJ
	case 'V':
		op = VJ
	default:
		return fmt.Errorf("unexpected character with code %d in table section: expected table character (S, D, C, V)", ap.peekByte())
	}

	ap.readByte()

	if ap.peekByte() != ':' {
		goto yield
	}

	ap.readByte()

	op |= ap.testFlag('I', IF)
	op |= ap.testFlag('E', EF)
	op |= ap.testFlag('M', MF)

	if ap.peekByte() != ':' {
		if b := ap.peekByte(); !isVoid(b) {
			return fmt.Errorf("unexpected character with code %d in flag section: expected flag character (I, E, M)", b)
		}
		goto yield
	}

	ap.readByte()

	op |= ap.testFlag('L', LC)
	op |= ap.testFlag('E', EC)
	op |= ap.testFlag('G', GC)

	if b := ap.peekByte(); !isVoid(b) {
		return fmt.Errorf("unexpected character with code %d in conditional section: expected conditional flag character (L, E, G)", b)
	}

yield:
	mac.code = append(mac.code, op)
	return nil
}

func (ap *AsmParser) testFlag(name byte, code Code) Code {
	if ap.peekByte() == name {
		ap.readByte()
		return code
	}

	return 0
}

func (ap *AsmParser) nextSect() {
	for c := ap.peekByte(); !isSpec(c); c = ap.peekByte() {
		ap.readByte()
	}
}

func isVoid(b byte) bool {
	return b == '\x00' || b == ' ' || b == '\n' || b == '\r' || b == '\t' || b == '\v'
}

func isSpec(b byte) bool {
	return b == '\x00' || b == ':' || b == '+' || b == '-'
}
