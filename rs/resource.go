package rs

import (
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
