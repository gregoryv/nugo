package graph

import (
	"fmt"
	"os"
)

// NewRootNode is a way to root a tree at a given point. Only
// difference from NewNode is it can contain / in the name.
func ExampleNewRootNode() {
	var (
		root = NewRootNode("/mnt/usb")
		a    = NewNode("a")
		b    = NewNode("b")
		c    = NewNode("file.txt")
	)
	root.Add(a, b)
	a.Add(c)
	root.Walk(LsRecursive(os.Stdout))
	// output:
	// /mnt/usb
	// /mnt/usb/a
	// /mnt/usb/a/file.txt
	// /mnt/usb/b
}

func ExampleNode_FirstChild_listAllChildren() {
	var (
		root = NewRootNode("/")
		a    = NewNode("a")
		b    = NewNode("b")
	)
	root.Add(a, b)
	c := root.FirstChild()
	for {
		if c == nil {
			break
		}
		fmt.Fprintln(os.Stdout, c.Name())
		c = c.sibling
	}
	// output:
	// a
	// b
}

func ExampleWalk() {
	var (
		root = NewRootNode("/")
		a    = NewNode("a")
		b    = NewNode("b")
		c    = NewNode("c")
	)
	root.Add(a, c)
	a.Add(b, NewNode("1"))
	c.Add(NewNode("x"), NewNode("y"))

	root.Walk(func(c *Node, abspath string, w *Walker) {
		fmt.Fprintln(os.Stdout, abspath)
		if abspath == "/c/x" {
			w.Stop()
		}
	})

	// output:
	// /
	// /a
	// /a/b
	// /a/1
	// /c
	// /c/x
}

func ExampleNode_DelChild() {
	root := NewRootNode("/")
	root.Make("etc", "bin", "tmp", "usr/")
	tmp := root.Find("/tmp")
	tmp.Make("y.txt", "dir")

	root.DelChild("etc")
	root.DelChild("no such")
	root.Find("/bin").DelChild("no such")
	tmp.DelChild("dir")
	tmp.DelChild("x.gz")
	root.Walk(LsRecursive(os.Stdout))
	// output:
	// /
	// /bin
	// /tmp
	// /tmp/y.txt
	// /usr%2F
}
