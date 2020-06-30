package rs

import "github.com/gregoryv/nugo"

type Operation uint

const (
	OpRead Operation = iota
	OpWrite
	OpExec
)

// Modes
func (me Operation) Modes() (n, u, g, o nugo.NodeMode) {
	switch me {
	case OpRead:
		return 04000, 00400, 00040, 00004
	case OpWrite:
		return 02000, 00200, 00020, 00002
	case OpExec:
		return 01000, 00100, 00010, 00001
	}
	panic("bad operation")
}

func (o Operation) String() string {
	switch o {
	case OpRead:
		return "read"
	case OpWrite:
		return "write"
	case OpExec:
		return "exec"
	default:
		return ""
	}
}
