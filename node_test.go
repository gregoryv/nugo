package rs

import (
	"fmt"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/fox"
)

func TestRootNode_Child(t *testing.T) {
	var (
		rn = NewRootNode("/", ModeDir)
	)
	if n := rn.Child(); n != nil {
		t.Fail()
	}
}

func ExampleNode_UnsetMode() {
	var (
		n = NewNode("node")
	)
	n.mode = ModeDir | 01755
	fmt.Println(n)
	n.UnsetMode(ModeDir)
	fmt.Println(n)
	// output:
	// d--xrwxr-xr-x 0 0 node
	// ---xrwxr-xr-x 0 0 node
}

func TestNode_SetResource(t *testing.T) {
	var (
		n = NewNode("val")
	)
	n.SetResource(1)
	if n.resource.(int) != 1 {
		t.Fail()
	}
}

func TestRootNode_Locate_itself(t *testing.T) {
	rn := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	rn.SetSeal(1, 1, 01755)
	rn.MakeAll("etc")
	ok, bad := asserter.NewMixed(t)

	ok(rn.Locate("/"))
	bad(rn.Locate("/ljlj"))
}

func TestRootNode_Locate(t *testing.T) {
	rn := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	bin := rn.Make("bin")
	rn.Make("etc")
	bin.Make("mkdir")

	b, err := rn.Locate("/bin/mkdir")
	if err != nil {
		t.Fatal(err)
	}
	got := b[len(b)-1].Name()
	if got != "mkdir" {
		rn.Walk(NodeLogger(t))
		t.Fail()
	}
}

func ExampleRootNode_Locate() {
	root := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	root.SetSeal(1, 1, 01755)
	root.MakeAll("etc", "tmp")

	nodes, _ := root.Locate("/etc")
	for _, node := range nodes {
		fmt.Println(node)
	}
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 etc
}

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

func Example_sortedDistinct() {
	root := NewRootNode("/", ModeSort|ModeDistinct)
	root.MakeAll("b", "a")
	root.Find("/b").MakeAll("2", "1", "1", "2")
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
	root.MakeAll("c", "b", "a")
	root.Find("/b").MakeAll("2", "1", "3", "0", "2.5")
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
	root.MakeAll("a", "b")
	root.Find("/mnt/usb/a").MakeAll("file.txt")
	root.Walk(NamePrinter(os.Stdout))
	// output:
	// /mnt/usb
	// /mnt/usb/a
	// /mnt/usb/a/file.txt
	// /mnt/usb/b
}

func ExampleNode_FirstChild_listAllChildren() {
	root := NewRoot("/")
	root.MakeAll("a", "b")
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
	root.MakeAll("a", "c")
	root.Find("/a").MakeAll("b", "1")
	root.Find("/c").MakeAll("x", "y")
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
	root.MakeAll("etc", "bin", "tmp", "usr/")
	tmp := root.Find("/tmp")
	tmp.MakeAll("y.txt", "dir")

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
