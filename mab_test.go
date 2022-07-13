package mabvm

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineRun(t *testing.T) {
	tests := []struct {
		name string
		init *Machine
		expt []Word
	}{
		{
			name: "incriment number",
			init: NewMachine(
				log.Default(),
				[]byte{
					VJ | EF,
				},
				[]Word{2, 3},
			),
			expt: []Word{2, 5},
		},
		{
			name: "decriment number",
			init: NewMachine(
				log.Default(),
				[]byte{
					VJ | IF | EF,
				},
				[]Word{2, 3},
			),
			expt: []Word{2, -1},
		},
		{
			name: "copy data",
			init: NewMachine(
				log.Default(),
				[]byte{
					DJ | IF, VJ, DJ | IF, VJ | IF,
				},
				[]Word{0, 7},
			),
			expt: []Word{7, 7},
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.init.Run()
			assert.Equal(t, test.expt, test.init.data, "it's so bad ...")
		})
	}
}
