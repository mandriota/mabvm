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
		switch err := ap.parseOpcodeOrNumber(mac); err {
		case nil:
		case io.EOF:
			return nil
		default:
			return err
		}
	}
}

func (ap *AsmParser) parseOpcodeOrNumber(mac *Machine) (err error) {
	ap.iterateSection()

	if ap.currentCharacter() != ':' {
		return ap.parseNumber(mac)
	}

	ap.iterateCharacter()

	op := Code(0)

	switch ap.currentCharacter() {
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

	ap.iterateCharacter()

	if cc := ap.currentCharacter(); cc != '\'' && cc != '"' {
		goto fini
	}

	ap.iterateCharacter()

	op |= ap.testFlag('I', IF)
	op |= ap.testFlag('E', EF)
	op |= ap.testFlag('M', MF)

	if cc := ap.currentCharacter(); cc != '"' {
		goto fini
	}

	ap.iterateCharacter()

	op |= ap.testFlag('L', LC)
	op |= ap.testFlag('E', EC)
	op |= ap.testFlag('G', GC)

fini:
	if cc := ap.currentCharacter(); !isVoid(cc) {
		return ap.buildError("character", "conditional flag")
	}

	mac.code = append(mac.code, op)
	return nil
}

func (ap *AsmParser) parseNumber(mac *Machine) (err error) {
	sign := int64(1)

	switch ap.currentCharacter() {
	case '+':
	case '-':
		sign = -1
	case '\x00':
		return io.EOF
	default:
		return ap.buildError("character", "sign")
	}

	ap.iterateCharacter()

	base := int64(0)

	switch ap.currentCharacter() {
	case 'b':
		base = 2
	case 'o':
		base = 8
	case 'd':
		base = 10
	case 'h':
		base = 16
	case '\x00':
		return ap.buildError("EOF", "base")
	default:
		return ap.buildError("character", "base")
	}

	ap.iterateCharacter()

	word := ap.parseNumberABS(base)

	if ap.currentCharacter() != '#' {
		mac.data = append(mac.data, word*sign)
	} else {
		ap.iterateCharacter()

		off := len(mac.data)

		mac.data = growSlice(mac.data, off+int(ap.parseNumberABS(base)))

		for i := range mac.data[off:] {
			mac.data[off+i] = word * sign
		}
	}

	return nil
}

func (ap *AsmParser) parseNumberABS(base int64) (word int64) {
	for {
		cc := ap.currentCharacter()

		dec := boolToByte(cc >= '0' && cc <= '9' && cc < '0'+byte(base))
		hex := boolToByte(cc >= 'A' && cc <= 'Z' && cc < '7'+byte(base))
		if dec+hex == 0 {
			break
		}

		word = word*base + int64(cc-'0'*dec-'7'*hex)

		ap.iterateCharacter()
	}

	return
}

func (ap *AsmParser) testFlag(name byte, code Code) Code {
	if ap.currentCharacter() == name {
		ap.iterateCharacter()
		return code
	}

	return 0
}

func (ap *AsmParser) iterateSection() {
	for cc := ap.currentCharacter(); !isSpec(cc); cc = ap.currentCharacter() {
		if cc == '\n' {
			ap.line++
		}

		ap.iterateCharacter()
	}
}

func (ap *AsmParser) iterateCharacter() {
	ap.pos++
}

func (ap *AsmParser) currentCharacter() byte {
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
