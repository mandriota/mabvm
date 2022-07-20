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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsmParser(t *testing.T) {
	tests := []struct {
		src string
		exp Opcode
	}{
		{src: ":S:IE:E", exp: SJ | IF | EF | EC},
	}

	ap := &AsmParser{}

	for _, test := range tests {
		ap.src = test.src
		ap.pos = 0

		v, err := ap.parseOpcode()
		if err != nil {
			t.Log(err)
			continue
		}
		assert.Equal(t, test.exp, v)
	}
}

func TestAsmParserMulti(t *testing.T) {
	test := `	:S:I
	:D:IE ; do something
	:V::E
	:C:I:LG
	`

	ops, err := (&AsmParser{
		src: test,
	}).Parse(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, op := range ops {
		t.Log(op)
	}
}

func BenchmarkAsmParserParseOpcode(b *testing.B) {
	ap := &AsmParser{
		src: "	:D:IEM:LEG",
	}

	for i := 0; i < b.N; i++ {
		ap.Reset(ap.src)
		ap.parseOpcode()
	}
}
