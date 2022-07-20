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
