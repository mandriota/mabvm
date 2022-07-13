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
