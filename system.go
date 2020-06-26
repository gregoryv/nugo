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

// ----------------------------------------
// syscalls
// ----------------------------------------
func (me *System) ls(w io.Writer) {
	me.rn.Walk(internal.NamePrinter(w))
}
