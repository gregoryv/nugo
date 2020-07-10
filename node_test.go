package nugo

import (
	"fmt"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestNode_AbsPath(t *testing.T) {
	rn := NewRootNode("/", ModeDir)
	a := rn.Make("a")
	b := a.Make("b")
	assert := asserter.New(t)
	assert().Equals(a.AbsPath(), "/a")
	assert().Equals(b.AbsPath(), "/a/b")
}

func TestNode_Copy(t *testing.T) {
	n := &Node{src: 1}
	c := n.Copy()
	n.src = 2
	if n.src == c.src {
		t.Fail()
	}
}

func TestNode_Source(t *testing.T) {
	n := &Node{src: 1}
	assert := asserter.New(t)
	assert(n.Source() != nil).Error("nil Source")
}

func TestNode_IsDir(t *testing.T) {
	file := &Node{}
	dir := &Node{mode: ModeDir}
	assert := asserter.New(t)
	assert(!file.IsDir())
	assert(dir.IsDir())
}

func TestNode_IsRoot(t *testing.T) {
	ok, bad := asserter.NewErrors(t)
	bad(NewNode("x").IsRoot()).Log("default")
	ok(NewRootNode("/x", 017555).IsRoot())
}

func TestRootNode_Parent(t *testing.T) {
	rn := NewRootNode("/", ModeDir)
	assert := asserter.New(t)
	assert(rn.Parent() == nil).Error("expect no child")
}

func TestRootNode_Child(t *testing.T) {
	rn := NewRootNode("/", ModeDir)
	assert := asserter.New(t)
	assert(rn.Child() == nil).Error("expect no child")
}

func TestRootNode_Sibling(t *testing.T) {
	rn := NewRootNode("/", ModeDir)
	assert := asserter.New(t)
	assert(rn.Sibling() == nil).Error("expect no sibling")
}

func ExampleNode_UnsetMode() {
	n := NewNode("node")
	n.mode = ModeDir | 01755
	fmt.Println(n)
	n.UnsetMode(ModeDir)
	fmt.Println(n)
	// output:
	// d--xrwxr-xr-x 0 0 node
	// ---xrwxr-xr-x 0 0 node
}

func TestNode_SetSource(t *testing.T) {
	n := NewNode("val")
	n.SetSource(1)
	if n.src.(int) != 1 {
		t.Fail()
	}
}

func TestRootNode_Find(t *testing.T) {
	rn := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	a := rn.Make("a")
	ok, bad := asserter.NewMixed(t)
	ok(rn.Find("/"))
	bad(rn.Find("/nosuch"))
	bad(a.Find("something")).Log("not root")
}

func Example_sortedDistinct() {
	root := NewRootNode("/", ModeSort|ModeDistinct)
	b := root.Make("b")
	b.MakeAll("2", "1", "1", "2")
	root.Make("a")
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
	b, _ := root.Find("/b")
	b.MakeAll("2", "1", "3", "0", "2.5")
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
	a := root.Make("a")
	a.MakeAll("file.txt")
	root.Make("b")
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

func ExampleRootNode_Walk() {
	root := NewRoot("/")
	a := root.Make("a")
	c := root.Make("c")
	a.MakeAll("b", "1")
	c.MakeAll("x", "y")
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
	tmp, _ := root.Find("/tmp")
	tmp.MakeAll("y.txt", "dir")

	root.DelChild("etc")
	root.DelChild("no such")
	bin, _ := root.Find("/bin")
	bin.DelChild("no such")
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
