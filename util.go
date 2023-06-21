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
	"sync"
	"unsafe"
)

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func growSlice[E any](s []E, l int) []E {
	if cap(s) < l {
		t := make([]E, l)
		copy(t, s)

		return t
	}

	return s[:l]
}

func byteOf[T any](v T) byte {
	return *(*byte)(unsafe.Pointer(&v))
}

func byteSliceOf[E any](s []E) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(s))),
		uintptr(len(s))*unsafe.Sizeof(*new(E)))
}

func synchronize(m *sync.RWMutex) {
	if m.TryLock() {
		m.Lock()
	} else {
		m.Unlock()
	}
}

func min[T integer](a, b T) T {
	if a < b {
		return a
	}

	return b
}

func max[T integer](a, b T) T {
	if a > b {
		return a
	}

	return b
}

func isVoid(b byte) bool {
	return b == '\x00' || b == ' ' || b == '\n' ||
		b == '\r' || b == '\t' || b == '\v'
}
