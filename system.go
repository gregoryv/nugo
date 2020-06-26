package graph

import (
	"io"

	"github.com/gregoryv/graph/internal"
)

func NewSystem() *System {
	rnMode := internal.ModeSort | internal.ModeDistinct
	rn := internal.NewRootNode("/", rnMode)
	rn.Make("bin")
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
	me.mounts("/").Walk(internal.NamePrinter(w))
}
