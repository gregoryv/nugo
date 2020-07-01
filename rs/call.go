package rs

import (
	"fmt"
	"path"

	"github.com/gregoryv/nugo"
)

type Syscall struct {
	*System
	acc *Account
}

// install resource at the absolute path
func (me *Syscall) Install(
	abspath string, resource interface{}, mode nugo.NodeMode,
) (
	*nugo.Node, error,
) {
	dir, name := path.Split(abspath)
	parent, err := me.stat(dir)
	if err != nil {
		return nil, err
	}
	if err := me.acc.permitted(OpWrite, parent.Seal()); err != nil {
		return nil, fmt.Errorf("Install: %v", err)
	}
	node := parent.Make(name)
	node.SetPerm(mode)
	node.SetResource(resource)
	node.UnsetMode(nugo.ModeDir)
	if resource != nil {

	}
	return node, nil
}

// ExecCmd creates and executes a new command with system defaults.
func (me *Syscall) ExecCmd(abspath string, args ...string) error {
	return me.Exec(NewCmd(abspath, args...))
}

// Exec executes the given command. Fails if e.g. resource is not
// Executable.
func (me *Syscall) Exec(cmd *Cmd) error {
	n, err := me.stat(cmd.Abspath)
	if err != nil {
		return err
	}
	switch r := n.Resource().(type) {
	case Executable:
		// If needed setuid can be checked and enforced here
		cmd.Sys = me
		return r.Exec(cmd)
	default:
		return fmt.Errorf("Cannot run %T", r)
	}
}

type Executable interface {
	Exec(*Cmd) error
}

// Mkdir creates the absolute path whith a given mode where the parent
// must exist.
func (me *Syscall) Mkdir(abspath string, mode nugo.NodeMode) (*ResInfo, error) {
	dir, name := path.Split(abspath)
	parent, err := me.stat(dir)
	if err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	if err := me.acc.permitted(OpWrite, parent.Seal()); err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	n := parent.Make(name)
	n.SetPerm(mode)
	return &ResInfo{node: n}, nil
}

// Stat returns the node of the abspath if account is allowed to reach
// it, ie. all nodes up to it must have execute flags set.
func (me *Syscall) Stat(abspath string) (*ResInfo, error) {
	n, err := me.stat(abspath)
	if err != nil {
		return nil, fmt.Errorf("Stat %v", err)
	}
	return &ResInfo{node: n}, nil
}

// stat returns the node of the abspath if account is allowed to reach
// it, ie. all nodes up to it must have execute flags set.
func (me *Syscall) stat(abspath string) (*nugo.Node, error) {
	rn := me.mounts(abspath)
	nodes, err := rn.Locate(abspath)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if err := me.acc.permitted(OpExec, n.Seal()); err != nil {
			return nil, fmt.Errorf("%s uid:%d: %v", abspath, me.acc.uid, err)
		}
	}
	return nodes[len(nodes)-1], nil
}
