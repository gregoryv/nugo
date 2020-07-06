package nugo

import (
	"fmt"
	"io"
	"path"

	"github.com/gregoryv/fox"
)

func NewWalker() *Walker {
	return &Walker{}
}

// Walker holds state of a walk.
type Walker struct {
	stopped bool
}

// Stop the Walker from your visitor.
func (me *Walker) Stop() { me.stopped = true }

// Walk calls the Visitor for the given node. The abspath is
// that of the child. Use empty string for root node.
func (w *Walker) Walk(parent, child *Node, abspath string, fn Visitor) {
	if child == nil || w.stopped {
		return
	}
	fn(parent, child, path.Join(abspath, child.Name()), w)
	w.Walk(child, child.child, path.Join(abspath, child.Name()), fn)
	w.Walk(parent, child.sibling, abspath, fn)
}

// Visitor is called during a walk with a specific node and the
// absolute path to that node. Use the given Walker to stop if needed.
// For root nodes the parent is nil.
type Visitor func(parent, child *Node, abspath string, w *Walker)

// ----------------------------------------

// NamePrinter writes abspath to the given writer.
func NamePrinter(w io.Writer) Visitor {
	return func(parent, child *Node, abspath string, Walker *Walker) {
		fmt.Fprintln(w, abspath)
	}
}

// NodePrinter writes permissions and ownership with each node
func NodePrinter(w io.Writer) Visitor {
	return func(parent, child *Node, abspath string, Walker *Walker) {
		fmt.Fprintln(w, child.Seal().String(), abspath)
	}
}

// NodeLogger logs permissions and ownership with each node
func NodeLogger(l fox.Logger) Visitor {
	return func(parent, child *Node, abspath string, Walker *Walker) {
		l.Log(fmt.Sprintf("%s %s", child.Seal().String(), abspath))
	}
}
