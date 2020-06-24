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

// NewNode returns a new node with the given name url path escaped.
func NewNode(name string) *Node {
	safe := url.PathEscape(name)
	return &Node{name: safe}
}

// Node names and links a sibling and a child.
type Node struct {
	name    string
	sibling *Node
	child   *Node
}

// Name returns the base name of a node
func (my *Node) Name() string { return my.name }

// Add adds each child in sequence.
func (me *Node) Add(children ...*Node) {
	for _, n := range children {
		me.addChild(n)
	}
}

// Make creates and adds the named children
func (me *Node) Make(names ...string) {
	for _, name := range names {
		n := NewNode(name)
		me.Add(n)
	}
}

// AppendChild adds the given node as the last child
func (me *Node) addChild(n *Node) {
	last := me.LastChild()
	if last == nil {
		me.child = n
		return
	}
	last.sibling = n
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
	for {
		sibling := next.sibling
		if sibling == nil {
			break
		}
		if sibling.name == name {
			next.sibling = next.sibling.sibling
			return sibling
		}
		next = sibling
	}
	return nil
}

// ----------------------------------------

// NewRootNode returns a new node with the name as is. It's the
// callers responsibility to make sure every basename is safe,
// Valid abspaths are "/" or "/mnt/usb"
func NewRootNode(abspath string) *RootNode {
	return &RootNode{
		Node: &Node{
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
	var visit func(*Node, string) *Node
	abspath = path.Clean(abspath)
	visit = func(n *Node, p string) *Node {
		if n == nil {
			return nil
		}
		prefix := path.Join(p, n.Name())
		if prefix == abspath {
			return n
		}
		if n := visit(n.child, prefix); n != nil {
			return n
		}
		return visit(n.sibling, p)
	}
	return visit(me.Node, "")
}

// Walk over each node until walker is stopped. Same as
//   NewWalker().Walk(root, "", fn)
func (me *RootNode) Walk(fn Visitor) {
	NewWalker().Walk(me.Node, "", fn)
}

// ----------------------------------------

func NewWalker() *Walker {
	return &Walker{}
}

// Walker holds state of a walk.
type Walker struct {
	stopped bool
}

// Stop the walker from your visitor.
func (me *Walker) Stop() { me.stopped = true }

// Walk calls the visitor for the given node. The abspath should be
// that of the parent. Use empty string for root graph.
func (w *Walker) Walk(n *Node, abspath string, fn Visitor) {
	if n == nil || w.stopped {
		return
	}
	fn(n, path.Join(abspath, n.Name()), w)
	w.Walk(n.child, path.Join(abspath, n.Name()), fn)
	w.Walk(n.sibling, abspath, fn)
}

// Visitor is called during a walk with a specific node and the
// absolute path to that node. Use the given walker to stop if needed.
type Visitor func(child *Node, abspath string, w *Walker)

// ----------------------------------------

// LsRecursive writes names of the children of n to w
func LsRecursive(w io.Writer) Visitor {
	return func(child *Node, abspath string, walker *Walker) {
		fmt.Fprintln(w, abspath)
	}
}
