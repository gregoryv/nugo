package sys

import (
	"fmt"
	"io"
	"path"

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
	sys.Install("/bin/mkdir", &mkdirCmd{}, Root, 00755)

	// Order is important until mkdir supports -p flag
	dirs := []string{
		"/etc",
		"/etc/accounts",
	}
	for _, dir := range dirs {
		sys.Exec(NewCmd("/bin/mkdir", "-m", "00755", dir), Root)
	}
	return sys
}

type System struct {
	rn *rs.RootNode
}

// Exec executes the given command as the account. Fails if
// e.g. resource is not Executable.
func (me *System) Exec(cmd *Cmd, acc *Account) error {
	n, err := me.Stat(cmd.Abspath, acc)
	if err != nil {
		return err
	}
	switch r := n.Resource().(type) {
	case Executable:
		// If needed setuid can be checked and enforced here
		cmd.Sys = &Syscall{acc: acc, sys: me}
		return r.Exec(cmd)
	default:
		return fmt.Errorf("Cannot run %T", r)
	}
}

type Executable interface {
	Exec(*Cmd) error
}

// mounts returns the mounting point of the abspath. Currently only
// "/" is available.
func (me *System) mounts(abspath string) *rs.RootNode {
	return me.rn
}

// ----------------------------------------
// syscalls
// ----------------------------------------

// dumprs writes the entire graph
func (me *System) dumprs(w io.Writer) {
	me.mounts("/").Walk(rs.NodePrinter(w))
}

// install resource at the absolute path
func (me *System) Install(
	abspath string, resource interface{}, acc *Account, mode rs.NodeMode,
) (
	*rs.Node, error,
) {
	dir, name := path.Split(abspath)
	n, err := me.Stat(dir, acc)
	if err != nil {
		return nil, err
	}
	if err := acc.Permitted(OpWrite, n.Seal()); err != nil {
		return nil, fmt.Errorf("Install: %v", err)
	}
	newNode := n.Make(name)
	newNode.SetPerm(mode)
	newNode.SetResource(resource)
	newNode.UnsetMode(rs.ModeDir)
	if resource != nil {

	}
	return newNode, nil
}

// Stat returns the node of the abspath if account is allowed to reach
// it, ie. all nodes up to it must have execute flags set.
func (me *System) Stat(abspath string, acc *Account) (*rs.Node, error) {
	rn := me.mounts(abspath)
	nodes, err := rn.Locate(abspath)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if err := acc.Permitted(OpExec, n.Seal()); err != nil {
			return nil, fmt.Errorf("Stat %s uid:%d: %v", abspath, acc.uid, err)
		}
	}
	return nodes[len(nodes)-1], nil
}
