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
	"errors"
	"io"
	"strconv"
)

type AsmParser struct {
	src string
	pos int

	line int
}

func NewAsmParser(src string) *AsmParser {
	return &AsmParser{src: src}
}

func (ap *AsmParser) Parse(mac *Machine) error {
	if mac == nil {
		return errors.New("machine should not be nil")
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
		return ap.buildError("character", "table")
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
			return ap.buildError("character", "control flag")
		}
		goto yield
	}

	ap.readByte()

	op |= ap.testFlag('L', LC)
	op |= ap.testFlag('E', EC)
	op |= ap.testFlag('G', GC)

	if b := ap.peekByte(); !isVoid(b) {
		return ap.buildError("character", "conditional flag")
	}

yield:
	mac.code = append(mac.code, op)
	return nil
}

func (ap *AsmParser) parseNumber(mac *Machine) (err error) {
	sign := int64(1)
	word := int64(0)
	base := int64(0)

	switch ap.peekByte() {
	case '+':
	case '-':
		sign = -1
	case '\x00':
		return io.EOF
	default:
		return ap.buildError("character", "sign")
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
		return ap.buildError("EOF", "base")
	default:
		return ap.buildError("character", "base")
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

func (ap *AsmParser) testFlag(name byte, code Code) Code {
	if ap.peekByte() == name {
		ap.readByte()
		return code
	}

	return 0
}

func (ap *AsmParser) nextSect() {
	for c := ap.peekByte(); !isSpec(c); c = ap.peekByte() {
		if c == '\n' {
			ap.line++
		}

		ap.readByte()
	}
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

func (ap *AsmParser) buildError(unexpect, expect string) error {
	return errors.New("line " + strconv.Itoa(ap.line) +
		": unexpected " + unexpect +
		": " + expect + " expected")
}

func isVoid(b byte) bool {
	return b == '\x00' || b == ' ' || b == '\n' || b == '\r' || b == '\t' || b == '\v'
}

func isSpec(b byte) bool {
	return b == '\x00' || b == ':' || b == '+' || b == '-'
}
