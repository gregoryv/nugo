package graph_test

import (
	"fmt"
	"os"

	"github.com/gregoryv/graph"
)

func Example() {
	root := graph.NewRoot("/")
	root.Make("etc", "tmp")
	root.Find("/tmp").Make("y.txt")
	root.Walk(graph.LsRecursive(os.Stdout))
	// output:
	// /
	// /etc
	// /tmp
	// /tmp/y.txt
}

func Example_graphManipulation() {
	root := graph.NewRootNode("/", graph.ModeSort)
	root.Make("etc", "tmp", "usr/")
	tmp := root.Find("/tmp")
	tmp.Make("y.txt", "dir")
	tmp.DelChild("dir")

	root.Walk(func(c *graph.Node, abspath string, w *graph.Walker) {
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
