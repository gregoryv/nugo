package sys

import (
	"fmt"
	"io"

	"github.com/gregoryv/rs"
)

func NewSystem() *System {
	rnMode := rs.ModeDir | rs.ModeSort | rs.ModeDistinct
	rn := rs.NewRootNode("/", rnMode)
	rn.SetSeal(1, 1, 01755)
	rn.Make("bin")

	etc := rn.Make("etc")
	etc.SetPerm(00755)
	etc.Make("accounts")

	return &System{
		rn: rn,
	}
}

type System struct {
	rn *rs.RootNode
}

// mounts returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) mounts(abspath string) *rs.RootNode {
	return me.rn
}

// ----------------------------------------
// syscalls
// ----------------------------------------

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.mounts("/").Walk(rs.NodePrinter(w))
}

// stat returns the node of the abspath if account is allowed to reach
// it, ie. all nodes up to it must have execute flags set.
func (me *System) stat(abspath string, acc *Account) (*rs.Node, error) {
	root := me.mounts(abspath)
	result := root.Locate(abspath)
	n := result.Node
	for {
		if n.Child() == nil {
			break
		}
		if err := acc.Permitted(OpExec, n.Seal()); err != nil {
			return nil, fmt.Errorf("stat %s uid:%d: %v", abspath, acc.uid, err)
		}
		n = n.Child()
	}
	return n, nil
}
