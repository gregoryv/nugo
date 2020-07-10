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
package nugo

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"
)

// Mutex for synchronizing all write and delete operations.
var mu sync.Mutex

// A NodeMode represents a node's mode and permission bits.
type NodeMode uint32

const (
	ModeDir      NodeMode = 1 << (32 - 1 - iota)
	ModeSort              // sorted by name
	ModeDistinct          // no duplicate children
	ModeRoot

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

// NewNode returns a new node with the given name url path escaped.
func NewNode(name string) *Node {
	safe := url.PathEscape(name)
	return &Node{name: safe}
}

// node names and links a sibling and a child.
type Node struct {
	name    string
	parent  *Node
	child   *Node
	sibling *Node

	uid  int
	gid  int
	mode NodeMode

	sync.RWMutex
	src interface{}
}

// AbsPath returns the absolute path to this node.
func (me *Node) AbsPath() string {
	if me.parent == nil {
		return me.name
	}
	return path.Join(me.parent.AbsPath(), me.name)
}

// Name returns the base name of a node
func (my *Node) Name() string { return my.name }

// Source returns the nodes resource or nil if none is set
func (my *Node) Source() interface{} { return my.src }

// Seal returns the access control seal of this node.
func (my *Node) Seal() Seal {
	return Seal{UID: my.uid, GID: my.gid, Mode: my.mode}
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

// SetSource of this node, use nil to clear
func (my *Node) SetSource(r interface{}) { my.src = r }

// IsDir returns true if ModeDir is set
func (me *Node) IsDir() bool { return me.mode&ModeDir != 0 }

// IsRoot
func (me *Node) IsRoot() error {
	if !me.isRoot() {
		return fmt.Errorf("%s not root", me.AbsPath())
	}
	return nil
}

func (me *Node) isRoot() bool { return me.mode&ModeRoot != 0 }

// UnsetMode todo
func (me *Node) UnsetMode(mask NodeMode) {
	me.mode = me.mode &^ mask
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
	mu.Lock()
	defer mu.Unlock()
	for _, n := range children {
		// inherit
		n.uid = me.uid
		n.gid = me.gid
		n.mode = me.mode &^ ModeRoot
		n.parent = me
		if n.mode&ModeDistinct == ModeDistinct {
			me.delChild(n.Name())
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

// Parent returns parent of this nodo.
func (my *Node) Parent() *Node { return my.parent }

// Child returns child of this node.
func (my *Node) Child() *Node { return my.child }

// Sibling
func (my *Node) Sibling() *Node { return my.sibling }

// Copy returns a copy of the node without relations and no source.
func (me *Node) Copy() *Node {
	cp := *me
	cp.child = nil
	cp.sibling = nil
	cp.src = nil
	return &cp
}

// DelChild removes the first child with the given name and returns the
// removed node
func (me *Node) DelChild(name string) *Node {
	mu.Lock()
	n := me.delChild(name)
	mu.Unlock()
	return n
}

func (me *Node) delChild(name string) *Node {
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

// String todo
func (me *Node) String() string {
	return fmt.Sprintf("%s %s", me.Seal(), me.name)
}

// ----------------------------------------

// NewRoot returns a new node with the name as is. It's the
// callers responsibility to make sure every basename is safe,
// Valid abspaths are "/" or "/mnt/usb".
// Defaults to mode ModeDir | ModeRoot.
func NewRoot(abspath string) *Node {
	return &Node{
		mode: ModeDir | ModeRoot,
		name: path.Clean(abspath),
	}
}

// NewRootNode returns a root node with the additional given mode.
func NewRootNode(abspath string, mode NodeMode) *Node {
	n := NewRoot(abspath)
	n.mode = n.mode | mode
	return n
}

// Find returns the node matching the absolute path. This node must be
// a root node.
func (me *Node) Find(abspath string) (*Node, error) {
	if err := me.IsRoot(); err != nil {
		return nil, fmt.Errorf("Find: %w", err)
	}
	fullname := path.Clean(abspath)
	var n *Node
	me.Walk(func(parent, child *Node, abspath string, w *Walker) {
		if fullname == abspath {
			n = child
			w.Stop()
		}
	})
	if n == nil {
		return nil, fmt.Errorf("%s no such directory or resource", abspath)
	}
	return n, nil
}

// Walk over each node until Walker is stopped.
func (me *Node) Walk(fn Visitor) {
	// todo adapt for walking from a child
	NewWalker().Walk(me.Parent(), me, "", fn)
}
