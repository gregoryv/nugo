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
		first:     true,
		recursive: true,
	}
}

// Walker holds state of a walk.
type Walker struct {
	first       bool
	recursive   bool
	skipChild   bool // do not enter a specific directory
	skipSibling bool
	skipVisit   bool
	stopped     bool
}

// Stop the Walker from your visitor.
func (me *Walker) Stop() { me.stopped = true }

// SetRecursive
func (me *Walker) SetRecursive(r bool) { me.recursive = r }

// Skip tells the walker not do descend into this node if it's in
// recursive mode and the node is a directory.
func (me *Walker) SkipChild() { me.skipChild = true }

// SkipSibling tells the walker to skip the next sibling.
// The field is reset after each node.
func (me *Walker) SkipSibling() { me.skipSibling = true }

// SkipVisit tells the walker not to call the visitor of the next
// node. The field is reset after each node.
func (me *Walker) SkipVisit() { me.skipVisit = true }

// Walk calls the Visitor for the given node. The abspath is
// that of the child. Use empty string for root node.
func (me *Walker) Walk(parent, child *Node, abspath string, fn Visitor) {
	if child == nil || me.stopped {
		return
	}
	if !me.first {
		me.skipChild = false
		me.skipSibling = false
		me.skipVisit = false
	}
	me.first = false
	// less allocation over child.AbsPath()
	childAbspath := path.Join(abspath, child.Name())
	if !me.skipVisit {
		fn(parent, child, childAbspath, me)
	}
	if (child.isRoot() || me.recursive) && !me.skipChild {
		me.Walk(child, child.child, childAbspath, fn)
	}
	if !me.skipSibling {
		me.Walk(parent, child.sibling, abspath, fn)
	}
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
