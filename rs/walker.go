package rs

import (
	"github.com/gregoryv/nugo"
)

func NewWalker(sys *Syscall) *Walker {
	return &Walker{
		w:   nugo.NewWalker(),
		sys: sys,
	}
}

type Walker struct {
	w   *nugo.Walker
	sys *Syscall
}

// SetRecursive
func (me *Walker) SetRecursive(r bool) { me.w.SetRecursive(r) }

// SkipSibling
func (me *Walker) SkipSibling() { me.w.SkipSibling() }
func (me *Walker) SkipVisit()   { me.w.SkipVisit() }

func (me *Walker) Walk(res *ResInfo, fn Visitor) error {
	// wrap the visitor with access control
	visitor := func(parent, child *nugo.Node, abspath string, w *nugo.Walker) {
		n := parent
		if parent == nil {
			n = child
		}
		if me.sys.acc.permitted(OpExec, n) != nil {
			w.SkipChild()
			return
		}
		var p *ResInfo
		if parent != nil {
			p = &ResInfo{parent}
		}
		c := &ResInfo{child}
		fn(p, c, abspath, w)
	}
	n := res.node
	parent := n.Parent()
	var abspath string
	if parent != nil {
		abspath = n.AbsPath()
	}
	me.w.Walk(parent, n, abspath, visitor)
	return nil
}
