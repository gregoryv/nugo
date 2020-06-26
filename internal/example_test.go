package internal

import (
	"fmt"
	"os"
)

func Example() {
	root := NewRoot("/")
	root.MakeAll("etc", "tmp")
	root.Find("/tmp").MakeAll("y.txt")
	root.Walk(NamePrinter(os.Stdout))
	// output:
	// /
	// /etc
	// /tmp
	// /tmp/y.txt
}

func Example_graphManipulation() {
	root := NewRootNode("/", ModeSort)
	root.MakeAll("etc", "tmp", "usr/")
	tmp := root.Find("/tmp")
	tmp.MakeAll("y.txt", "dir")
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
