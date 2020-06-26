package internal

import (
	"fmt"
	"os"
)

func ExampleNodePrinter() {
	root := NewRootNode("/", ModeSort|ModeDistinct)
	root.Walk(NodePrinter(os.Stdout))
	// output:
	// ------------ 0 0 /
}

func Example_sortedDistinct() {
	root := NewRootNode("/", ModeSort|ModeDistinct)
	root.Make("b", "a")
	root.Find("/b").Make("2", "1", "1", "2")
	root.Walk(NamePrinter(os.Stdout))
	// output:
	// /
	// /a
	// /b
	// /b/1
	// /b/2
}

func Example_sorted() {
	root := NewRootNode("/", ModeSort)
	root.Make("c", "b", "a")
	root.Find("/b").Make("2", "1", "3", "0", "2.5")
	root.Walk(NamePrinter(os.Stdout))
	// output:
	// /
	// /a
	// /b
	// /b/0
	// /b/1
	// /b/2
	// /b/2.5
	// /b/3
	// /c
}

// NewRootNode is a way to root a tree at a given point. Only
// difference from NewNode is it can contain / in the name.
func ExampleNewRootNode() {
	root := NewRoot("/mnt/usb")
	root.Make("a", "b")
	root.Find("/mnt/usb/a").Make("file.txt")
	root.Walk(NamePrinter(os.Stdout))
	// output:
	// /mnt/usb
	// /mnt/usb/a
	// /mnt/usb/a/file.txt
	// /mnt/usb/b
}

func ExampleNode_FirstChild_listAllChildren() {
	root := NewRoot("/")
	root.Make("a", "b")
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
	root := NewRoot("/")
	root.Make("a", "c")
	root.Find("/a").Make("b", "1")
	root.Find("/c").Make("x", "y")
	root.Walk(func(parent, c *Node, abspath string, w *Walker) {
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
	root := NewRoot("/")
	root.Make("etc", "bin", "tmp", "usr/")
	tmp := root.Find("/tmp")
	tmp.Make("y.txt", "dir")

	root.DelChild("etc")
	root.DelChild("no such")
	root.Find("/bin").DelChild("no such")
	tmp.DelChild("dir")
	tmp.DelChild("x.gz")
	root.Walk(NamePrinter(os.Stdout))
	// output:
	// /
	// /bin
	// /tmp
	// /tmp/y.txt
	// /usr%2F
}
