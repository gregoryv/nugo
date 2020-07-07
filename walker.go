package nugo

import (
	"fmt"
	"io"
	"path"

	"github.com/gregoryv/fox"
)

// NewWalker returns a recursive walker
func NewWalker() *Walker {
	return &Walker{
		recursive: true,
	}
}

// Walker holds state of a walk.
type Walker struct {
	recursive bool
	stopped bool
}

// Stop the Walker from your visitor.
func (me *Walker) Stop() { me.stopped = true }

// SetRecursive
func (me *Walker) SetRecursive(r bool)  { me.recursive = r }


// Walk calls the Visitor for the given node. The abspath is
// that of the child. Use empty string for root node.
func (me *Walker) Walk(parent, child *Node, abspath string, fn Visitor) {
	if child == nil || me.stopped {
		return
	}
	fn(parent, child, path.Join(abspath, child.Name()), me)
	if parent == nil || me.recursive {
		me.Walk(child, child.child, path.Join(abspath, child.Name()), fn)
	}
	me.Walk(parent, child.sibling, abspath, fn)
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
