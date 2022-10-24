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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsmParser_parseOpcode(t *testing.T) {
	tests := []struct {
		name   string
		source string
		expect []Code
	}{
		{
			name:   "maximal case",
			source: ":D:IEM:LEG",
			expect: []Code{DJ | IF | EF | MF | LC | EC | GC},
		},
		{
			name:   "minimal case",
			source: ":V",
			expect: []Code{VJ},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := &Machine{}

			ap := NewAsmParser(test.source)

			assert.Equal(t, nil, ap.parseOpcode(mc))
			assert.Equal(t, test.expect, mc.code)
		})
	}
}

func TestAsmParser_parseNumber(t *testing.T) {
	tests := []struct {
		name   string
		source string
		expect []Word
	}{
		{
			name:   "one binary digit",
			source: "+b1",
			expect: []Word{0b1},
		},
		{
			name:   "one octave digit",
			source: "+o5",
			expect: []Word{05},
		},
		{
			name:   "one decimal digit",
			source: "+d9",
			expect: []Word{9},
		},
		{
			name:   "one hexadecimal digit",
			source: "+hB",
			expect: []Word{0xB},
		},
		{
			name:   "positive binary number",
			source: "+b101010010001111011011",
			expect: []Word{0b101010010001111011011},
		},
		{
			name:   "positive octave number",
			source: "+o7654325034562567",
			expect: []Word{07654325034562567},
		},
		{
			name:   "positive decimal number",
			source: "+d4993509343295043294",
			expect: []Word{4993509343295043294},
		},
		{
			name:   "positive hexadecimal number",
			source: "+h391AFE38F9C2DBAA",
			expect: []Word{0x391AFE38F9C2DBAA},
		},
		{
			name:   "negative binary number",
			source: "-b1011101010010010010101011110",
			expect: []Word{-0b1011101010010010010101011110},
		},
		{
			name:   "negative octave number",
			source: "-o76523251643042042043",
			expect: []Word{-076523251643042042043},
		},
		{
			name:   "negative decimal number",
			source: "-d3928419499493694382",
			expect: []Word{-3928419499493694382},
		},
		{
			name:   "negative hexadecimal number",
			source: "-hAF483A34FE49BC2",
			expect: []Word{-0xAF483A34FE49BC2},
		},
		{
			name:   "one binary digit duplicate",
			source: "+b1#100",
			expect: []Word{0b1, 0b1, 0b1, 0b1},
		},
		{
			name:   "one octave digit duplicate",
			source: "+o7#5",
			expect: []Word{07, 07, 07, 07, 07},
		},
		{
			name:   "one decimal digit duplicate",
			source: "+d9#8",
			expect: []Word{9, 9, 9, 9, 9, 9, 9, 9},
		},
		{
			name:   "one hexadecimal digit duplicate",
			source: "+hA#B",
			expect: []Word{0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA, 0xA},
		},
		{
			name:   "binary number duplicate",
			source: "+b1101#111",
			expect: []Word{0b1101, 0b1101, 0b1101, 0b1101, 0b1101, 0b1101, 0b1101},
		},
		{
			name:   "octave number duplicate",
			source: "+o756#10",
			expect: []Word{0756, 0756, 0756, 0756, 0756, 0756, 0756, 0756},
		},
		{
			name:   "decimal number duplicate",
			source: "+d783287#11",
			expect: []Word{783287, 783287, 783287, 783287, 783287, 783287, 783287, 783287, 783287, 783287, 783287},
		},
		{
			name:   "hexadecimal number duplicate",
			source: "-h0000A000#1",
			expect: []Word{-0xA000},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := &Machine{}

			ap := NewAsmParser(test.source)

			assert.Equal(t, nil, ap.parseNumber(mc))
			assert.Equal(t, test.expect, mc.data)
		})
	}
}

func TestAsmParserParse(t *testing.T) {
	test := struct {
		source string
		expect *Machine
	}{
		source: `	+b1010 +d643 +o746
	:V:E
	:D
	+b10111
	:S:I +hC0FFEE#6
	:V:I`,
		expect: &Machine{
			code: []Code{VJ | EF, DJ, SJ | IF, VJ | IF},
			data: []Word{0b1010, 643, 0746, 0b10111, 0xC0FFEE, 0xC0FFEE, 0xC0FFEE, 0xC0FFEE, 0xC0FFEE, 0xC0FFEE},
		},
	}

	ap := NewAsmParser(test.source)

	mc := &Machine{}

	assert.Equal(t, nil, ap.Parse(mc))
	assert.Equal(t, test.expect, mc)
}

func BenchmarkAsmParserParse(b *testing.B) {
	src := `+d484932984 :D:I :V :S :D:I :V:I :S :D`

	ap := NewAsmParser(src)
	mc := &Machine{}

	for i := 0; i < b.N; i++ {
		ap.src = src
		ap.pos = 0
		ap.Parse(mc)
	}
}
