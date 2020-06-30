package rs

import "github.com/gregoryv/nugo"

// ResInfo describes a resource and is returned by Stat
type ResInfo struct {
	node *nugo.Node
}

// Name returns the name of the file
func (me *ResInfo) Name() string { return me.node.Name() }
