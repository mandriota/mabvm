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

func TestMachineRun(t *testing.T) {
	tests := []struct {
		name string
		init func() *Machine
		expt []Word
	}{
		{
			name: "incriment number",
			init: func() *Machine {
				m := NewMachine(
					[]Code{VJ | EF},
					[]Word{2, 3},
					nil,
				)
				return m
			},
			expt: []Word{5, 3},
		},
		{
			name: "decriment number",
			init: func() *Machine {
				m := NewMachine(
					[]Code{VJ | IF | EF},
					[]Word{3, 2},
					nil,
				)
				return m
			},
			expt: []Word{1, 2},
		},
		{
			name: "copy data",
			init: func() *Machine {
				m := NewMachine(
					[]Code{VJ | EF},
					[]Word{0, 7},
					nil,
				)
				return m
			},
			expt: []Word{7, 7},
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mac := test.init()
			mac.Show()
			assert.Equal(t, test.expt, mac.data, "it's so bad ...")
		})
	}
}

func BenchmarkMachineRun(b *testing.B) {
	mac := NewMachine(
		[]Code{VJ, DJ | IF, SJ, VJ, DJ | IF, SJ, VJ, DJ | IF, SJ, VJ, DJ | IF, SJ},
		[]Word{0, 123},
		nil,
	)

	for i := 0; i < b.N; i++ {
		mac.Show()
		mac.codP = 0
	}
}
