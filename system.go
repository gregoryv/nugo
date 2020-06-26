package nugo

import (
	"io"

	"github.com/gregoryv/nugo/internal"
)

func NewSystem() *System {
	rnMode := internal.ModeDir | internal.ModeSort | internal.ModeDistinct
	rn := internal.NewRootNode("/", rnMode)
	rn.SetSeal(1, 1, 01755)
	n := rn.Make("bin")
	n.SetSeal(1, 1, 01755)
	return &System{
		rn: rn,
	}
}

type System struct {
	rn *internal.RootNode
}

// mounts returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) mounts(abspath string) *internal.RootNode {
	return me.rn
}

// ----------------------------------------
// syscalls
// ----------------------------------------

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.mounts("/").Walk(internal.NodePrinter(w))
}
