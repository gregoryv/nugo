package rs

import (
	"bytes"
	"fmt"

	"github.com/gregoryv/nugo"
)

// ResInfo describes a resource and is returned by Stat
type ResInfo struct {
	node *nugo.Node
}

// Name returns the name of the file
func (me *ResInfo) Name() string { return me.node.Name() }

// IsDir returns nil if the resource is a directory
func (me *ResInfo) IsDir() error {
	if !me.node.IsDir() {
		return fmt.Errorf("IsDir: %s not a directory", me.node.Name())
	}
	return nil
}

func newResource(n *nugo.Node, op operation) *Resource {
	return &Resource{
		node: n,
		op:   op,
	}
}

// Resource wraps access to the underlying node
type Resource struct {
	node   *nugo.Node
	op     operation
	unlock func()
	buf    *bytes.Buffer // used for reading node source
}

// SetSource sets the src of the underlying node. Returns error if it's readonly.
func (me *Resource) SetSource(src interface{}) error {
	if me.op&OpWrite == 0 {
		return fmt.Errorf("SetSource: %s read only", me.node.Name())
	}
	switch src := src.(type) {
	case string:
	case []byte:
	case Executable:
	default:
		return fmt.Errorf("SetSource: %T cannot be used as source", src)
	}
	me.node.SetSource(src)
	return nil
}

// Read
func (me *Resource) Read(b []byte) (int, error) {
	if me.op&OpRead == 0 {
		return 0, fmt.Errorf("Read: %s not readable", me.node.Name())
	}
	if me.buf == nil {
		return 0, fmt.Errorf("Read: unreadable source")
	}
	return me.buf.Read(b)
}

// Write
func (me *Resource) Write(p []byte) (int, error) {
	if me.op&OpWrite == 0 {
		return 0, fmt.Errorf("Write: %s not writable", me.node.Name())
	}
	return me.buf.Write(p)
}

// Close
func (me *Resource) Close() error {
	me.buf = nil
	me.unlock()
	return nil
}
