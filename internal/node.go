/*
Package internal provides a directed internal implementation.

  root
    |
    +-- child
        sibling
        sibling
        |    |
        *    +-- child
                 sibling
                 ...

*/
package internal

import (
	"fmt"
	"io"
	"net/url"
	"path"
)

// A NodeMode represents a node's mode and permission bits.
type NodeMode uint32

const (
	ModeSort     NodeMode = 1 << (32 - 1 - iota) // sorted by name
	ModeDistinct                                 // no duplicate children
)

// newNode returns a new node with the given name url path escaped.
func NewNode(name string) *Node {
	safe := url.PathEscape(name)
	return &Node{name: safe}
}

// node names and links a sibling and a child.
type Node struct {
	name    string
	mode    NodeMode
	sibling *Node
	child   *Node

	resource interface{}
}

// Name returns the base name of a node
func (my *Node) Name() string { return my.name }

// Make creates and adds the named children
func (me *Node) Make(names ...string) {
	for _, name := range names {
		me.Add(NewNode(name))
	}
}

// Add adds each child in sequence according to the NodeMode of the
// parent node.
func (me *Node) Add(children ...*Node) {
	for _, n := range children {
		n.mode = me.mode
		if n.mode&ModeDistinct == ModeDistinct {
			me.DelChild(n.Name())
		}
		switch {
		case n.mode&ModeSort == ModeSort:
			me.insert(n)
		default:
			me.append(n)
		}
	}
}

func (me *Node) append(n *Node) {
	last := me.LastChild()
	if last == nil {
		me.child = n
		return
	}
	last.sibling = n
}

// insert the node sorted by name
func (me *Node) insert(n *Node) {
	switch {
	case me.child == nil:
		me.child = n
	case n.Name() < me.child.Name():
		n.sibling = me.child
		me.child = n
	default:
		me.insertSibling(me.child, n)
	}
}

// insertSibling inserts n as a sibling of c
func (me *Node) insertSibling(c, n *Node) {
	for {
		if c.sibling == nil {
			c.sibling = n
			return
		}
		if n.Name() < c.sibling.Name() {
			n.sibling = c.sibling
			c.sibling = n
			return
		}
		c = c.sibling
	}
}

// FirstChild returns the first child or nil if there are no children.
func (me *Node) FirstChild() *Node { return me.child }

// LastChild returns the last child or nil if there are no children.
func (me *Node) LastChild() *Node {
	if me.child == nil {
		return nil
	}
	last := me.child
	for {
		if last.sibling == nil {
			break
		}
		last = last.sibling
	}
	return last
}

// DelChild removes the first child with the given name and returns the
// removed node
func (me *Node) DelChild(name string) *Node {
	if me.child == nil {
		return nil
	}
	next := me.child
	if next.name == name {
		me.child = next.sibling
		return next
	}
	return me.delSibling(me.child, name)
}

func (me *Node) delSibling(c *Node, name string) *Node {
	for {
		sibling := c.sibling
		if sibling == nil {
			break
		}
		if sibling.name == name {
			c.sibling = c.sibling.sibling
			return sibling
		}
		c = sibling
	}
	return nil
}

// ----------------------------------------

// NewRoot returns a rootNode with no special mode set.
func NewRoot(abspath string) *RootNode {
	return newRootNode(abspath, 0)
}

// NewRootNode returns a new node with the name as is. It's the
// callers responsibility to make sure every basename is safe,
// Valid abspaths are "/" or "/mnt/usb"
func newRootNode(abspath string, mode NodeMode) *RootNode {
	return &RootNode{
		Node: &Node{
			mode: mode,
			name: path.Clean(abspath),
		},
	}
}

type RootNode struct {
	*Node
}

// Find returns the node matching the absolute path starting at the
// root.
func (me *RootNode) Find(abspath string) *Node {
	fullname := path.Clean(abspath)
	var n *Node
	me.Walk(func(parent, child *Node, abspath string, w *Walker) {
		if fullname == abspath {
			n = child
			w.Stop()
		}
	})
	return n
}

// Walk over each node until Walker is stopped. Same as
//   NewWalker().Walk(root, "", fn)
func (me *RootNode) Walk(fn Visitor) {
	newWalker().Walk(nil, me.Node, "", fn)
}

// ----------------------------------------

func newWalker() *Walker {
	return &Walker{}
}

// Walker holds state of a walk.
type Walker struct {
	stopped bool
}

// Stop the Walker from your visitor.
func (me *Walker) Stop() { me.stopped = true }

// Walk calls the Visitor for the given node. The abspath should be
// that of the parent. Use empty string for root internal.
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
