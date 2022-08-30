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
		[]Code{DJ | IF, SJ | IF | EF, VJ | EF, VJ},
		append(append([]Word{0, textW}, make([]Word, BlockSize*2-3)...), 8188),
		nil,
	)

	buf := bytes.NewBuffer(nil)
	wrt := NewWriter(buf, mac.data[4096:])
	go wrt.Show()

	mac.Bind(&wrt.Mutex, wrt.Blocks())
	mac.Show()

	// wait for writer
	runtime.GOMAXPROCS(runtime.NumGoroutine())
	runtime.Gosched()

	assert.Equal(
		t,
		append(make([]byte, 1<<15-16), textB[:]...),
		buf.Bytes(),
	)
}
