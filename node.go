/*

Package rs provides a resource system with unix style access control.

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

*/
package rs

import (
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
)

// A NodeMode represents a node's mode and permission bits.
type NodeMode uint32

const (
	ModeDir      NodeMode = 1 << (32 - 1 - iota)
	ModeSort              // sorted by name
	ModeDistinct          // no duplicate children

	ModeType NodeMode = ModeSort | ModeDistinct
	ModePerm NodeMode = 07777
)

// String returns the bits as string using rwx notation for each bit.
func (m NodeMode) String() string {
	var s strings.Builder
	s.WriteByte(mChar(m, ModeDir, 'd'))
	s.WriteByte(mChar(m, 04000, 'r'))
	s.WriteByte(mChar(m, 02000, 'w'))
	s.WriteByte(mChar(m, 01000, 'x'))
	s.WriteByte(mChar(m, 00400, 'r'))
	s.WriteByte(mChar(m, 00200, 'w'))
	s.WriteByte(mChar(m, 00100, 'x'))
	s.WriteByte(mChar(m, 00040, 'r'))
	s.WriteByte(mChar(m, 00020, 'w'))
	s.WriteByte(mChar(m, 00010, 'x'))
	s.WriteByte(mChar(m, 00004, 'r'))
	s.WriteByte(mChar(m, 00002, 'w'))
	s.WriteByte(mChar(m, 00001, 'x'))
	return s.String()
}

func mChar(m, mask NodeMode, c byte) byte {
	if m&mask == mask {
		return c
	}
	return '-'
}

// newNode returns a new node with the given name url path escaped.
func NewNode(name string) *Node {
	safe := url.PathEscape(name)
	return &Node{name: safe}
}

// node names and links a sibling and a child.
type Node struct {
	name    string
	sibling *Node
	child   *Node

	uid  int
	gid  int
	mode NodeMode

	resource interface{}
}

// Name returns the base name of a node
func (my *Node) Name() string { return my.name }

// Seal returns the access control seal of this node.
func (my *Node) Seal() *Seal {
	return &Seal{uid: my.uid, gid: my.gid, perm: my.mode}
}

// SetUID sets the owner id of this node.
func (my *Node) SetUID(uid int) { my.uid = uid }

// SetGID sets group owner of this node.
func (my *Node) SetGID(gid int) { my.gid = gid }

// SetPerm permission bits of this node.
func (my *Node) SetPerm(perm NodeMode) {
	mode := my.mode &^ ModePerm // clear previous
	my.mode = mode | perm
}

// SetSeal sets ownership and permission mode of the this node.
func (my *Node) SetSeal(uid, gid int, mode NodeMode) {
	my.SetUID(uid)
	my.SetGID(gid)
	my.SetPerm(mode)
}

// Make creates and adds the named child returning the new node.
// The new node as ModeDir set.
func (me *Node) Make(name string) *Node {
	n := NewNode(name)
	n.SetSeal(me.uid, me.gid, me.mode) // inherit parent
	me.Add(n)
	return n
}

// MakeAll creates and adds the named children
func (me *Node) MakeAll(names ...string) {
	for _, name := range names {
		me.Make(name)
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

// copy returns a copy of the given node without child, sibling
// relations or resource
func (me *Node) copy() *Node {
	return &Node{
		name: me.name,
		uid:  me.uid,
		gid:  me.gid,
		mode: me.mode,
	}
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

// NewRoot returns a rootNode with ModeDir set.
func NewRoot(abspath string) *RootNode {
	return NewRootNode(abspath, ModeDir)
}

// NewRootNode returns a new node with the name as is. It's the
// callers responsibility to make sure every basename is safe,
// Valid abspaths are "/" or "/mnt/usb"
func NewRootNode(abspath string, mode NodeMode) *RootNode {
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

// Locate returns a new root node with each child set.
func (me *RootNode) Locate(abspath string) *RootNode {
	fullname := path.Clean(abspath)
	newRoot := NewRootNode(me.name, me.mode)
	newRoot.Node = me.Node.copy()

	n := newRoot.Node
	me.Walk(func(parent, child *Node, abspath string, w *Walker) {
		if parent == nil { // skip parent
			return
		}
		if strings.Index(fullname, abspath) == 0 {
			c := child.copy()
			n.child = c
			if fullname == abspath {
				w.Stop()
			}
			n = c
		}
	})
	return newRoot
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
// that of the parent. Use empty string for root rs.
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
