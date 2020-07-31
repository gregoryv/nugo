/*
Package nugo provides a graph with unix style modes.

The graph is a set of linked nodes

  root
    |
    +-- child
        sibling
        sibling
        |    |
        *    +-- child
                 sibling
                 ...

Each node references its parent aswell.
In addition to the standard unix permission rwxrwxrwx another set of
rwx are added to indicate anonymous access control. Permission bits
for "Other" means other authenticated.

             rwxrwxrwxrwx    ModePerm
              n  u  g  o
              |  |  |  |
  aNonymous --+  |  |  |
       User -----+  |  |
      Group --------+  |
      Other -----------+

Operations on the graph are not synchronized, this is left to any
system using it.

*/
package nugo

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
)

// NewRootNode returns a root node with the additional given mode.
func NewRootNode(abspath string, mode NodeMode) *Node {
	n := NewRoot(abspath)
	n.Mode = n.Mode | mode
	return n
}

// NewRoot returns a new node with the name as is. It's the
// callers responsibility to make sure every basename is safe,
// Valid abspaths are "/" or "/mnt/usb".
// Defaults to mode ModeDir | ModeRoot.
func NewRoot(abspath string) *Node {
	n := &Node{
		Name: path.Clean(abspath),
	}
	n.Mode = ModeDir | ModeRoot
	return n
}

// NewNode returns a new node with the given name url path escaped.
func NewNode(name string) *Node {
	safe := url.PathEscape(name)
	return &Node{Name: safe}
}

// Node names and links a sibling and a child.
type Node struct {
	Name string
	Seal
	Content interface{}

	sync.RWMutex
	Parent  *Node
	Child   *Node
	sibling *Node
}

// AbsPath returns the absolute path to this node.
func (me *Node) AbsPath() string {
	if me.Parent == nil {
		return me.Name
	}
	return path.Join(me.Parent.AbsPath(), me.Name)
}

// SetPerm permission bits of this node.
func (my *Node) SetPerm(perm NodeMode) {
	Mode := my.Mode &^ ModePerm // clear previous
	my.Mode = Mode | perm
}

// CheckDir returns nil if ModeDir is set
func (me *Node) CheckDir() error {
	if !me.IsDir() {
		return fmt.Errorf("%s not directory", me.AbsPath())
	}
	return nil
}

// IsDir returns true if ModeDir is set.
func (me *Node) IsDir() bool { return me.Mode&ModeDir != 0 }

// CheckRoot returns nil if ModeRoot is set
func (me *Node) CheckRoot() error {
	if !me.IsRoot() {
		return fmt.Errorf("%s not root", me.AbsPath())
	}
	return nil
}

// IsRoot returns true if ModeRoot is set.
func (me *Node) IsRoot() bool { return me.Mode&ModeRoot != 0 }

// SetMode sets mode of this node.
func (my *Node) SetMode(Mode NodeMode) {
	my.Mode = my.Mode | Mode
}

// UnsetMode unsets the given mode
func (me *Node) UnsetMode(mask NodeMode) {
	me.Mode = me.Mode &^ mask
}

// SetSeal sets ownership and permission mode of the this node.
func (my *Node) SetSeal(uid, gid int, Mode NodeMode) {
	my.UID = uid
	my.GID = gid
	my.SetPerm(Mode)
}

// Make creates and adds the named child returning the new node.
// See Add method for mode inheritence.
func (me *Node) Make(name string) *Node {
	n := NewNode(name)
	me.Add(n)
	return n
}

// MakeAll creates and adds the named children
func (me *Node) MakeAll(names ...string) []*Node {
	nodes := make([]*Node, len(names))
	for i, name := range names {
		nodes[i] = me.Make(name)
	}
	return nodes
}

// Add adds each child in sequence according to the NodeMode of the
// parent node. Add blocks if another add is in progress.
// For ModeDistinct an existing node with the same name is replaced.
func (me *Node) Add(children ...*Node) {
	for _, n := range children {
		// inherit
		n.UID = me.UID
		n.GID = me.GID
		n.Mode = me.Mode &^ ModeRoot
		n.Parent = me
		if n.Mode&ModeDistinct == ModeDistinct {
			me.delChild(n.Name)
		}
		switch {
		case n.Mode&ModeSort == ModeSort:
			me.insert(n)
		default:
			me.append(n)
		}
	}
}

func (me *Node) append(n *Node) {
	last := me.LastChild()
	if last == nil {
		me.Child = n
		return
	}
	last.sibling = n
}

// insert the node sorted by name
func (me *Node) insert(n *Node) {
	switch {
	case me.Child == nil:
		me.Child = n
	case n.Name < me.Child.Name:
		n.sibling = me.Child
		me.Child = n
	default:
		me.insertSibling(me.Child, n)
	}
}

// insertSibling inserts n as a sibling of c
func (me *Node) insertSibling(c, n *Node) {
	for {
		if c.sibling == nil {
			c.sibling = n
			return
		}
		if n.Name < c.sibling.Name {
			n.sibling = c.sibling
			c.sibling = n
			return
		}
		c = c.sibling
	}
}

// FirstChild returns the first child or nil if there are no children.
func (me *Node) FirstChild() *Node { return me.Child }

// LastChild returns the last child or nil if there are no children.
func (me *Node) LastChild() *Node {
	if me.Child == nil {
		return nil
	}
	last := me.Child
	for {
		if last.sibling == nil {
			break
		}
		last = last.sibling
	}
	return last
}

// Sibling
func (my *Node) Sibling() *Node { return my.sibling }

// Copy returns a copy of the node without relations and no source.
func (me *Node) Copy() *Node {
	cp := *me
	cp.Child = nil
	cp.sibling = nil
	cp.Content = nil
	return &cp
}

// DelChild removes the first child with the given name and returns the
// removed node
func (me *Node) DelChild(name string) *Node {
	n := me.delChild(name)
	return n
}

func (me *Node) delChild(name string) *Node {
	if me.Child == nil {
		return nil
	}
	next := me.Child
	if next.Name == name {
		me.Child = next.sibling
		return next
	}
	return me.delSibling(me.Child, name)
}

func (me *Node) delSibling(c *Node, name string) *Node {
	for {
		sibling := c.sibling
		if sibling == nil {
			break
		}
		if sibling.Name == name {
			c.sibling = c.sibling.sibling
			return sibling
		}
		c = sibling
	}
	return nil
}

// String todo
func (me *Node) String() string {
	return fmt.Sprintf("%s %s", me.Seal, me.Name)
}

// Find returns the node matching the absolute path. This node must be
// a root node.
func (me *Node) Find(abspath string) (*Node, error) {
	if err := me.CheckRoot(); err != nil {
		return nil, fmt.Errorf("Find: %w", err)
	}
	fullname := path.Clean(abspath)
	if fullname == me.Name {
		return me, nil
	}
	var n *Node
	visitor := func(Child *Node, abspath string, w *Walker) {
		if fullname == abspath {
			n = Child
			w.Stop()
		}
		if !strings.HasPrefix(fullname, abspath) {
			w.SkipChild() // don't descend into this child
		}
	}
	walker := NewWalker()
	walker.Walk(me, visitor)
	if n == nil {
		return nil, fmt.Errorf("%s no such directory or resource", abspath)
	}
	return n, nil
}

// Walk over each node until Walker is stopped.
func (me *Node) Walk(fn Visitor) {
	NewWalker().Walk(me, fn)
}
