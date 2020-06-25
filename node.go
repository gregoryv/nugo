/*
Package graph provides a directed graph implementation.

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
package graph

import (
	"fmt"
	"io"
	"net/url"
	"path"
)

// A nodeMode represents a node's mode and permission bits.
type nodeMode uint32

const (
	ModeSort     nodeMode = 1 << (32 - 1 - iota) // sorted by name
	ModeDistinct                                 // no duplicate children
)

// newNode returns a new node with the given name url path escaped.
func newnode(name string) *node {
	safe := url.PathEscape(name)
	return &node{name: safe}
}

// node names and links a sibling and a child.
type node struct {
	name    string
	mode    nodeMode
	sibling *node
	child   *node

	resource interface{}
}

// Name returns the base name of a node
func (my *node) Name() string { return my.name }

// Make creates and adds the named children
func (me *node) Make(names ...string) {
	for _, name := range names {
		me.Add(newnode(name))
	}
}

// Add adds each child in sequence according to the nodeMode of the
// parent node.
func (me *node) Add(children ...*node) {
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

func (me *node) append(n *node) {
	last := me.LastChild()
	if last == nil {
		me.child = n
		return
	}
	last.sibling = n
}

// insert the node sorted by name
func (me *node) insert(n *node) {
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
func (me *node) insertSibling(c, n *node) {
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
func (me *node) FirstChild() *node { return me.child }

// LastChild returns the last child or nil if there are no children.
func (me *node) LastChild() *node {
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
func (me *node) DelChild(name string) *node {
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

func (me *node) delSibling(c *node, name string) *node {
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
func NewRoot(abspath string) *rootNode {
	return newRootNode(abspath, 0)
}

// NewRootNode returns a new node with the name as is. It's the
// callers responsibility to make sure every basename is safe,
// Valid abspaths are "/" or "/mnt/usb"
func newRootNode(abspath string, mode nodeMode) *rootNode {
	return &rootNode{
		node: &node{
			mode: mode,
			name: path.Clean(abspath),
		},
	}
}

type rootNode struct {
	*node
}

// Find returns the node matching the absolute path starting at the
// root.
func (me *rootNode) Find(abspath string) *node {
	fullname := path.Clean(abspath)
	var n *node
	me.Walk(func(parent, child *node, abspath string, w *walker) {
		if fullname == abspath {
			n = child
			w.Stop()
		}
	})
	return n
}

// Walk over each node until walker is stopped. Same as
//   NewWalker().Walk(root, "", fn)
func (me *rootNode) Walk(fn visitor) {
	newWalker().Walk(nil, me.node, "", fn)
}

// ----------------------------------------

func newWalker() *walker {
	return &walker{}
}

// walker holds state of a walk.
type walker struct {
	stopped bool
}

// Stop the walker from your visitor.
func (me *walker) Stop() { me.stopped = true }

// Walk calls the visitor for the given node. The abspath should be
// that of the parent. Use empty string for root graph.
func (w *walker) Walk(parent, child *node, abspath string, fn visitor) {
	if child == nil || w.stopped {
		return
	}
	fn(parent, child, path.Join(abspath, child.Name()), w)
	w.Walk(child, child.child, path.Join(abspath, child.Name()), fn)
	w.Walk(parent, child.sibling, abspath, fn)
}

// visitor is called during a walk with a specific node and the
// absolute path to that node. Use the given walker to stop if needed.
// For root nodes the parent is nil.
type visitor func(parent, child *node, abspath string, w *walker)

// ----------------------------------------

// NamePrinter writes abspath to the given writer.
func NamePrinter(w io.Writer) visitor {
	return func(parent, child *node, abspath string, walker *walker) {
		fmt.Fprintln(w, abspath)
	}
}
