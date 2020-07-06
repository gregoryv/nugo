package nugo

import (
	"os"

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
	root.Make("etc") // inherits parent mode
	l := fox.NewSyncLog(os.Stdout)
	root.Walk(NodeLogger(l))
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /etc
}
