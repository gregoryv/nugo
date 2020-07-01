package rs

import (
	"fmt"
	"sync"

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

type Resource struct {
	readOnly bool
	mu       sync.Mutex
	node     *nugo.Node
}

// SetSource sets the src of the underlying node. Returns error if it's readonly.
func (me *Resource) SetSource(src Src) error {
	if me.readOnly {
		return fmt.Errorf("SetSource: %s read only", me.node.Name())
	}

	// maybe limit to certain types here
	me.mu.Lock()
	me.node.SetSource(src)
	me.mu.Unlock()
	return nil
}
