package nugo

import (
	"io"

	"github.com/gregoryv/rs"
)

// NewSystem returns a system with a default resources resembling a
// unix filesystem.
func NewSystem() *System {
	rnMode := rs.ModeDir | rs.ModeSort | rs.ModeDistinct
	rn := rs.NewRootNode("/", rnMode)
	rn.SetSeal(1, 1, 01755)
	rn.Make("bin")

	sys := &System{rn: rn}
	syscall := &Syscall{System: sys, acc: Root}
	syscall.Install("/bin/mkdir", &mkdirCmd{}, 00755)

	// Order is important until mkdir supports -p flag
	dirs := []string{
		"/etc",
		"/etc/accounts",
	}
	for _, dir := range dirs {
		syscall.Exec(NewCmd("/bin/mkdir", "-m", "00755", dir))
	}
	return sys
}

type System struct {
	rn *rs.RootNode
}

// mounts returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) mounts(abspath string) *rs.RootNode {
	return me.rn
}

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.mounts("/").Walk(rs.NodePrinter(w))
}
