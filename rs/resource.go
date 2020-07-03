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

// Resource wraps access to the underlying node
type Resource struct {
	readOnly bool
	node     *nugo.Node
	unlock   func()
	buf      *bytes.Buffer // used for reading node source
}

// SetSource sets the src of the underlying node. Returns error if it's readonly.
func (me *Resource) SetSource(src interface{}) error {
	if me.readOnly {
		return fmt.Errorf("SetSource: %s read only", me.node.Name())
	}
	// maybe limit to certain types here
	me.node.SetSource(src)
	return nil
}

// Read
func (me *Resource) Read(b []byte) (int, error) {
	if me.buf == nil {
		src := me.node.Source()
		switch src := src.(type) {
		case []byte:
			me.buf = bytes.NewBuffer(src)
		case string:
			me.buf = bytes.NewBufferString(src)
		}
	}
	return me.buf.Read(b)
}

// Close
func (me *Resource) Close() error {
	me.buf = nil
	me.unlock()
	return nil
}
