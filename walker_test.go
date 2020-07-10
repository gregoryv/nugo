package nugo

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/fox"
)

func ExampleNodePrinter() {
	root := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	root.SetSeal(1, 1, 01755)
	root.Make("etc") // inherits parent mode
	root.Walk(NodePrinter(os.Stdout))
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /etc
}

func ExampleNodeLogger() {
	root := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	root.SetSeal(1, 1, 01755)
	tmp := root.Make("tmp")
	tmp.Make("sub")
	l := fox.NewSyncLog(os.Stdout)
	root.Walk(NodeLogger(l))
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /tmp
	// d--xrwxr-xr-x 1 1 /tmp/sub
}

func TestWalker_Walk_recursive(t *testing.T) {
	root := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	root.SetSeal(1, 1, 01755)
	tmp := root.Make("tmp")
	tmp.Make("sub")
	walker := NewWalker()
	walker.SetRecursive(false) // recursive by default
	var buf bytes.Buffer
	visitor := func(p, c *Node, abspath string, w *Walker) {
		fmt.Fprintln(&buf, abspath)
	}
	walker.Walk(nil, root.Node, "", visitor)
	if strings.Contains(buf.String(), "sub") {
		t.Error("contains sub:", buf.String())
	}
}

func TestWalker_Skip(t *testing.T) {
	root := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	root.SetSeal(1, 1, 01755)
	tmp := root.Make("tmp")
	tmp.Make("sub")
	walker := NewWalker()
	var buf bytes.Buffer
	visitor := func(p, c *Node, abspath string, w *Walker) {
		fmt.Fprintln(&buf, abspath)
		if abspath == "/tmp" {
			w.Skip()
		}
	}
	walker.Walk(nil, root.Node, "", visitor)
	if strings.Contains(buf.String(), "sub") {
		t.Error("contains sub:", buf.String())
	}
}
