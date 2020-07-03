/*
Package rs provides a resource system which enforces unix style access control.
*/
package rs

import (
	"io"

	"github.com/gregoryv/nugo"
)

// NewSystem returns a system with a default resources resembling a
// unix filesystem.
func NewSystem() *System {
	rnMode := nugo.ModeDir | nugo.ModeSort | nugo.ModeDistinct
	rn := nugo.NewRootNode("/", rnMode)
	rn.SetSeal(1, 1, 01755)
	rn.Make("bin")

	sys := &System{
		rn: rn,
	}
	syscall := &Syscall{System: sys, acc: Root}
	syscall.Install("/bin/mkdir", &mkdirCmd{}, 00755)

	// Order is important until mkdir supports -p flag
	dirs := []struct {
		abspath string
		mode    string
	}{
		{"/etc", "00755"},
		{"/etc/accounts", "00755"},
		{"/tmp", "07777"},
	}
	for _, d := range dirs {
		syscall.Exec(NewCmd("/bin/mkdir", "-m", d.mode, d.abspath))
	}
	return sys
}

type System struct {
	rn *nugo.RootNode
}

// Use returns a syscall for the given account
func (me *System) Use(acc *Account) *Syscall {
	return &Syscall{System: me, acc: acc}
}

// mounts returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) mounts(abspath string) *nugo.RootNode {
	return me.rn
}

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.mounts("/").Walk(nugo.NodePrinter(w))
}
