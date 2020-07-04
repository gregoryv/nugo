package rs

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path"

	"github.com/gregoryv/nugo"
)

type Syscall struct {
	*System // probably should hide this
	acc     *Account
}

// RemoveAll
func (me *Syscall) RemoveAll(abspath string) error {
	rn := me.mounts(abspath)
	nodes, err := rn.Locate(abspath)
	if err != nil {
		return wrap("RemoveAll", err)
	}
	for _, n := range nodes {
		if err := me.acc.permitted(OpExec, n.Seal()); err != nil {
			return fmt.Errorf("%s uid:%d: %v", abspath, me.acc.uid, err)
		}
	}
	last := len(nodes) - 1
	nodes[last-1].DelChild(nodes[last].Name())
	return nil
}

// Open resource for reading. Underlying source must be string or []byte.
// If resource is open for writing this call blocks.
func (me *Syscall) Open(abspath string) (*Resource, error) {
	n, err := me.stat(abspath)
	if err != nil {
		return nil, fmt.Errorf("Open: %s", err)
	}
	r := newResource(n, OpRead)
	r.unlock = n.RUnlock
	src := n.Source()
	switch src := src.(type) {
	case []byte:
		r.buf = bytes.NewBuffer(src)
	default:
		// todo figure out how to read Any source
		return nil, fmt.Errorf("Open: %s(%T) non readable source", abspath, src)
	}
	// Resource must be closed to unlock
	n.RLock()
	return r, nil
}

// Create returns a new resource for writing. Fails if existing
// resource is directory. Caller must close resource.
func (me *Syscall) Create(abspath string) (*Resource, error) {
	rif, _ := me.Stat(abspath)
	if rif != nil && rif.IsDir() == nil {
		return nil, fmt.Errorf("Create: %s is a directory", abspath)
	}
	dir, name := path.Split(abspath)
	parent, err := me.Stat(dir)
	if err != nil {
		return nil, wrap("Create", err)
	}
	if err := me.acc.permitted(OpWrite, parent.node.Seal()); err != nil {
		return nil, wrap("Create", err)
	}
	n := parent.node.Make(name)
	n.SetPerm(00644)
	n.UnsetMode(nugo.ModeDir)
	n.Lock()
	r := newResource(n, OpWrite)
	r.buf = &bytes.Buffer{}
	r.unlock = n.Unlock
	return r, nil
}

// SaveAs save src to the given abspath. Fails if abspath already exists.
func (me *Syscall) SaveAs(abspath string, src interface{}) error {
	if _, err := me.Stat(abspath); err == nil {
		return fmt.Errorf("SaveAs: %s exists", abspath)
	}
	w, err := me.Create(abspath)
	if err != nil {
		return wrap("SaveAs", err)
	}
	defer w.Close()
	return wrap("SaveAs", gob.NewEncoder(w).Encode(src))
}

// Save save src to the given abspath. Overwrites existing resource.
func (me *Syscall) Save(abspath string, src interface{}) error {
	rif, _ := me.Stat(abspath)
	if rif != nil && rif.IsDir() == nil {
		return fmt.Errorf("Save: %s is directory", abspath)
	}
	w, err := me.Create(abspath)
	if err != nil {
		return wrap("Save", err)
	}
	defer w.Close()
	return wrap("Save", gob.NewEncoder(w).Encode(src))
}

// Load
func (me *Syscall) Load(res interface{}, abspath string) error {
	r, err := me.Open(abspath)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	return wrap("Load", gob.NewDecoder(r).Decode(res))
}

// Install resource at the absolute path
func (me *Syscall) Install(abspath string, cmd Executable, mode nugo.NodeMode,
) (*ResInfo, error) {
	dir, name := path.Split(abspath)
	parent, err := me.Stat(dir)
	if err != nil {
		return nil, wrap("Install", err)
	}
	if err := me.acc.permitted(OpWrite, parent.node.Seal()); err != nil {
		return nil, wrap("Install", err)
	}
	n := parent.node.Make(name)
	n.SetPerm(mode)
	n.SetSource(cmd)
	n.UnsetMode(nugo.ModeDir)
	return &ResInfo{node: n}, nil
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
	switch src := n.Source().(type) {
	case Executable:
		// If needed setuid can be checked and enforced here
		cmd.Sys = me
		return src.Exec(cmd)
	default:
		return fmt.Errorf("Cannot run %T", src)
	}
}

type Executable interface {
	Exec(*Cmd) error
}

// mkdir creates the absolute path whith a given mode where the parent
// must exist.
func (me *Syscall) mkdir(abspath string, mode nugo.NodeMode) (*ResInfo, error) {
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

func wrap(prefix string, err error) error {
	if err != nil {
		return fmt.Errorf("%s: %w", prefix, err)
	}
	return nil
}

// mount creates a root node for the given path.
func (me *Syscall) mount(abspath string, mode nugo.NodeMode) error {
	rn := nugo.NewRootNode(abspath, mode)
	rn.SetSeal(me.acc.uid, me.acc.gid(), 01755)
	me.System.rn = rn
	return nil
}
