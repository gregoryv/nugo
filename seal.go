package nugo

import "fmt"

type Sealed interface {
	Seal() Seal
}

type Seal struct {
	UID  int // owner
	GID  int // group
	Mode NodeMode
}

func (l Seal) String() string {
	return fmt.Sprintf("%s %d %d", l.Mode.String(), l.UID, l.GID)
}
