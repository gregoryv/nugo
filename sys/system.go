package sys

import (
	"flag"
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

	sys := &System{
		rn: rn,
	}
	sys.Install("/bin/mkdir", &MkdirCmd{}, Root, 00755)

	// todo use mkdir to create subsequent directories
	// Root.Exec("/bin/mkdir", &Mkdir{
	//    Abspath: "/etc",
	//    Mode: 00755,
	// })
	etc := rn.Make("etc")
	etc.SetPerm(00755)
	etc.Make("accounts")

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

// Mkdir creates the director given by abspath
func (me *System) Mkdir(abspath string, mode rs.NodeMode, acc *Account) (*rs.Node, error) {
	dir, name := path.Split(abspath)
	n, err := me.Stat(dir, acc)
	if err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	if err := acc.Permitted(OpWrite, n.Seal()); err != nil {
		return nil, fmt.Errorf("Mkdir: %w", err)
	}
	newNode := n.Make(name)
	newNode.SetPerm(mode)
	return newNode, nil
}

type MkdirCmd struct{}

func (me *MkdirCmd) Exec(c *Command, args ...string) error {
	flags := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	mode := flags.Uint("m", 00755, "mode for new directory")
	if err := flags.Parse(args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	_, err := c.sys.Mkdir(abspath, rs.NodeMode(*mode), c.acc)
	return err
}

type Command struct {
	sys *System
	acc *Account
}
