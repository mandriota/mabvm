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

const (
	IF = 1 << iota << 2 // inv. flag
	EF                  // ext. flag : y :
	MF                  // mut. flag
	LC                  // low. cnd. flag
	EC                  // equ. cnd. flag
	GC                  // gre. cnd. flag
)

const (
	SJ = iota // jump src head/tail table :: z
	DJ        // jump dst head/tail table :: z
	CJ        // jump code head/tail table :: z
	VJ        // jump value head/tail table : x : z
)

const JMask = SJ | DJ | CJ | VJ
