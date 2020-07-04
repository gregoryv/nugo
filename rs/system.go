/*
Package rs provides a resource system which enforces unix style access control.

Resources are stored as nugo.Nodes and can either have a []byte slice
as source or implement the Executable interface. Using the Save and
Load syscalls, structs are gob encoded and decoded to an access
controlled resource.

Anonymous account has uid,gid 0,0 whereas the Root account 1,1.

*/
package rs

import (
	"io"

	"github.com/gregoryv/nugo"
)

// NewSystem returns a system with installed resources resembling a
// unix filesystem.
func NewSystem() *System {
	sys := &System{}
	asRoot := Root.Use(sys)
	asRoot.mount("/", nugo.ModeDir|nugo.ModeSort|nugo.ModeDistinct)
	installSys(sys)
	return sys
}

// installSys creates default resources on the system. Should only be
// called once on one system.
func installSys(sys *System) {
	asRoot := Root.Use(sys)
	asRoot.mkdir("/bin", 01755)
	asRoot.mkdir("/etc", 00755)
	asRoot.mkdir("/etc/accounts", 00755)
	asRoot.mkdir("/tmp", 07777)
	asRoot.Install("/bin/mkdir", &mkdirCmd{}, 00755)
}

type System struct {
	rn *nugo.RootNode
}

// Use returns a syscall for the given account
func (me *System) Use(acc *Account) *Syscall {
	return &Syscall{System: me, acc: acc}
}

// rootNode returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) rootNode(abspath string) *nugo.RootNode {
	return me.rn
}

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.rootNode("/").Walk(nugo.NodePrinter(w))
}
