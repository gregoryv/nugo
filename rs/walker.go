package rs

import (
	"path"

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

func (me *Walker) Walk(abspath string, fn Visitor) error {
	n, err := me.sys.stat(abspath)
	if err != nil {
		return err
	}
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
	me.w.Walk(n.Parent(), n, path.Dir(abspath), visitor)
	return nil
}
