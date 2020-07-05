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
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/gregoryv/fox"
	"github.com/gregoryv/nugo"
)

// NewSystem returns a system with installed resources resembling a
// unix filesystem.
func NewSystem() *System {
	sys := &System{
		mounts: make(map[string]*nugo.RootNode),
	}
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
	asRoot.Install("/bin/ls", &lsCmd{}, 00755)
}

type System struct {
	mounts  map[string]*nugo.RootNode
	auditer fox.Logger // Used audit Syscall.Exec calls
}

// SetAuditer sets the auditer for Syscall.Exec calls
func (me *System) SetAuditer(auditer fox.Logger) *System {
	me.auditer = auditer
	return me
}

func (me *System) mount(rn *nugo.RootNode) error {
	abspath := path.Clean(rn.Name())
	if _, found := me.mounts[abspath]; found {
		return fmt.Errorf("mount: %s already exists", abspath)
	}
	me.mounts[abspath] = rn
	return nil
}

// rootNode returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) rootNode(abspath string) *nugo.RootNode {
	rn := me.mounts["/"]
	for p, n := range me.mounts {
		if strings.Index(abspath, p) == 0 {
			if len(n.Name()) > len(rn.Name()) {
				rn = n
			}
		}
	}
	return rn
}

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.rootNode("/").Walk(nugo.NodePrinter(w))
}
