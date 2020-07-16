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

// Walk calls the Visitor for the given node and all related.
func (me *Walker) Walk(node *Node, fn Visitor) {
	me.walk(node, "", fn)
}

func (me *Walker) walk(node *Node, parent string, fn Visitor) {
	if node == nil || me.stopped {
		return
	}
	if !me.first {
		me.skipChild = false
		me.skipSibling = false
		me.skipVisit = false
	}
	me.first = false
	// less allocation over node.AbsPath()
	abspath := path.Join(parent, node.Name())
	if !me.skipVisit {
		fn(node, abspath, me)
	}
	if (node.isRoot() || me.recursive) && !me.skipChild {
		me.walk(node.child, abspath, fn)
	}
	if !me.skipSibling {
		me.walk(node.sibling, parent, fn)
	}
}

// Visitor is called during a walk with a specific node and the
// absolute path to that node. Use the given Walker to stop if needed.
// For root nodes the parent is nil.
type Visitor func(node *Node, abspath string, w *Walker)

// ----------------------------------------

// NamePrinter writes abspath to the given writer.
func NamePrinter(writer io.Writer) Visitor {
	return func(child *Node, abspath string, w *Walker) {
		fmt.Fprintln(writer, abspath)
	}
}

// NodePrinter writes permissions and ownership with each node
func NodePrinter(writer io.Writer) Visitor {
	return func(child *Node, abspath string, w *Walker) {
		fmt.Fprintln(writer, child.Seal().String(), abspath)
	}
}

// NodeLogger logs permissions and ownership with each node
func NodeLogger(l fox.Logger) Visitor {
	return func(child *Node, abspath string, w *Walker) {
		l.Log(fmt.Sprintf("%s %s", child.Seal().String(), abspath))
	}
}
