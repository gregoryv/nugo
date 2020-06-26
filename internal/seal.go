package internal

import "fmt"

type Sealed interface {
	Seal() Seal
}

type Seal struct {
	uid  int // owner
	gid  int // group
	perm NodeMode
}

func (l Seal) String() string {
	return fmt.Sprintf("%s %d %d", l.perm.String(), l.uid, l.gid)
}
