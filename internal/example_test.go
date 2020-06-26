package internal

import (
	"fmt"
	"os"
)

func Example() {
	root := NewRoot("/")
	root.Make("etc", "tmp")
	root.Find("/tmp").Make("y.txt")
	root.Walk(NamePrinter(os.Stdout))
	// output:
	// /
	// /etc
	// /tmp
	// /tmp/y.txt
}

func Example_graphManipulation() {
	root := newRootNode("/", ModeSort)
	root.Make("etc", "tmp", "usr/")
	tmp := root.Find("/tmp")
	tmp.Make("y.txt", "dir")
	tmp.DelChild("dir")

	root.Walk(func(parent, child *Node, abspath string, w *Walker) {
		fmt.Fprintln(os.Stdout, abspath)
		if abspath == "/tmp/y.txt" {
			w.Stop()
		}
	})
	// output:
	// /
	// /etc
	// /tmp
	// /tmp/y.txt
}
