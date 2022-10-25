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

// MabVM - Stack Virtual Machine named after Mab - Queen of The Faires.
//
// Its memory designed as a linear array divided in 8x4KiB blocks where each block has own mutex.
// When some subject write to not self memory, there must be locked mutex.
//
// Each its instruction is derived from 1x 2-bit sequence jump code and 6x 1-bit flags.
// Each instruction increases own sequence number by value which defaultly is equal to 1.
package mabvm

const (
	// IF - Inversion Flag. Inverts value.
	IF = 1 << iota << 2

	// EF - Extension Flag. Extends value from source.
	EF

	// MF - Mutex Flag. Blocks execution until the current block is changed.
	MF

	// LC - Lower Conditional flag. Will executed only if source is lower than destination.
	LC

	// EC - Equal Conditional flag. Will executed only if source is equal to destination.
	EC

	// GC - Greater Conditional Flag. Will executed only if source is greater than destination.
	GC
)

const (
	// SJ - Source Jump. Increments source pointer.
	SJ = iota

	// DJ - Destination Jump. Increments destination pointer.
	DJ

	// CJ - Code Jump. Increments code pointer.
	CJ

	// VJ - Value Jump. Increments value.
	VJ
)

// JMask - mask for sequence jumps.
const JMask = SJ | DJ | CJ | VJ
