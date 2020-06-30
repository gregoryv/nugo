package nugo

import (
	"fmt"
	"path"

	"github.com/gregoryv/rs"
)

type Syscall struct {
	*System
	acc *Account
}

// install resource at the absolute path
func (me *Syscall) Install(
	abspath string, resource interface{}, mode rs.NodeMode,
) (
	*rs.Node, error,
) {
	dir, name := path.Split(abspath)
	n, err := me.Stat(dir)
	if err != nil {
		return nil, err
	}
	if err := me.acc.Permitted(OpWrite, n.Seal()); err != nil {
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

// Exec executes the given command as the account. Fails if
// e.g. resource is not Executable.
func (me *Syscall) Exec(cmd *Cmd) error {
	n, err := me.Stat(cmd.Abspath)
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
func (me *Syscall) Mkdir(abspath string, mode rs.NodeMode) (*rs.Node, error) {
	dir, name := path.Split(abspath)
	parent, err := me.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	if err := me.acc.Permitted(OpWrite, parent.Seal()); err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	node := parent.Make(name)
	node.SetPerm(mode)
	return node, nil
}

// Stat returns the node of the abspath if account is allowed to reach
// it, ie. all nodes up to it must have execute flags set.
func (me *Syscall) Stat(abspath string) (*rs.Node, error) {
	rn := me.mounts(abspath)
	nodes, err := rn.Locate(abspath)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if err := me.acc.Permitted(OpExec, n.Seal()); err != nil {
			return nil, fmt.Errorf("Stat %s uid:%d: %v", abspath, me.acc.uid, err)
		}
	}
	return nodes[len(nodes)-1], nil
}
