package sys

import (
	"fmt"
	"path"

	"github.com/gregoryv/rs"
)

type Syscall struct {
	acc *Account
	sys *System
}

// Mkdir creates the absolute path whith a given mode where the parent
// must exist.
func (me *Syscall) Mkdir(abspath string, mode rs.NodeMode) (
	*rs.Node, error,
) {
	dir, name := path.Split(abspath)
	parent, err := me.sys.Stat(dir, me.acc)
	if err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	if err := me.acc.Permitted(OpWrite, parent.Seal()); err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	node := parent.Make(name)
	node.SetPerm(mode)
	return node, nil
}
