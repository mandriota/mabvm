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

import "unsafe"

func growSlice[E any](s []E, l int) []E {
	if cap(s) < l {
		t := make([]E, l)
		copy(t, s)

		return t
	}

	return s[:l]
}

func byteSliceOf[E any](s []E) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&s[0])),
		uintptr(len(s))*unsafe.Sizeof(s[0]))
}

func byteOf[T any](v T) byte {
	return *(*byte)(unsafe.Pointer(&v))
}

func isVoid(b byte) bool {
	return b == '\x00' || b == ' ' || b == '\n' ||
		b == '\r' || b == '\t' || b == '\v'
}

func isSpec(b byte) bool {
	return b == '\x00' || b == ':' || b == '+' || b == '-'
}
