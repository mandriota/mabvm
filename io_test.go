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
	"bytes"
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestWriterRun(t *testing.T) {
	textB := [8]byte{'H', 'i', ',', ' ', 'M', 'A', 'B', '\n'}
	textW := *(*Word)(unsafe.Pointer(&textB))

	mac := NewMachine(
		[]Code{VJ, VJ, DJ | EF, VJ},
		append(make([]Word, 8187), textW, 0, 4093, 0, 8186),
		new(MutexTab),
	)

	buf := bytes.NewBuffer(nil)
	wrt := NewWriter(buf, mac.data[:4096], mac.data)

	go wrt.Show()

	mac.Bind(&wrt.Mutex, wrt.Blocks())
	mac.Show()

	for buf.Len() == 0 {
		runtime.Gosched()
	}

	assert.Equal(
		t,
		textB[:],
		buf.Bytes(),
	)
}
