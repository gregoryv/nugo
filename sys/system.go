package sys

import (
	"io"

	"github.com/gregoryv/nugo"
)

func NewSystem() *System {
	rnMode := nugo.ModeDir | nugo.ModeSort | nugo.ModeDistinct
	rn := nugo.NewRootNode("/", rnMode)
	rn.SetSeal(1, 1, 01755)
	n := rn.Make("bin")
	n.SetSeal(1, 1, 01755)
	return &System{
		rn: rn,
	}
}

type System struct {
	rn *nugo.RootNode
}

// mounts returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) mounts(abspath string) *nugo.RootNode {
	return me.rn
}

// ----------------------------------------
// syscalls
// ----------------------------------------

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.mounts("/").Walk(nugo.NodePrinter(w))
}
