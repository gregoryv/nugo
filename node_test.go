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
	n := &Node{Content: 1}
	c := n.Copy()
	n.Content = 2
	if n.Content == c.Content {
		t.Fail()
	}
}

func TestNode_CheckDir(t *testing.T) {
	file := &Node{}
	dir := &Node{Seal: Seal{Mode: ModeDir}}
	ok, bad := asserter.NewErrors(t)
	ok(dir.CheckDir())
	bad(file.CheckDir())
}

func TestNode_IsDir(t *testing.T) {
	file := &Node{}
	dir := &Node{Seal: Seal{Mode: ModeDir}}
	assert := asserter.New(t)
	assert(!file.IsDir())
	assert(dir.IsDir())
}

func TestNode_CheckRoot(t *testing.T) {
	ok, bad := asserter.NewErrors(t)
	bad(NewNode("x").CheckRoot()).Log("default")
	ok(NewRootNode("/x", 017555).CheckRoot())
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
	n.Mode = ModeDir | 01755
	fmt.Println(n)
	n.UnsetMode(ModeDir)
	fmt.Println(n)
	// output:
	// d--xrwxr-xr-x 0 0 node
	// ---xrwxr-xr-x 0 0 node
}

func ExampleNode_SetMode() {
	n := NewNode("node")
	n.Mode = 01755
	fmt.Println(n)
	n.SetMode(ModeDir)
	fmt.Println(n)
	// output:
	// ---xrwxr-xr-x 0 0 node
	// d--xrwxr-xr-x 0 0 node
}

func TestNode_Find(t *testing.T) {
	rn := NewRootNode("/", ModeDir|ModeSort|ModeDistinct)
	a := rn.Make("a")
	ok, bad := asserter.NewMixed(t)
	ok(rn.Find("/a"))
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
	root.Walk(func(c *Node, abspath string, w *Walker) {
		fmt.Fprintln(os.Stdout, abspath)
		if abspath == "/c/x" {
			w.Stop()
		}
	})
	// output:
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
	// /bin
	// /tmp
	// /tmp/y.txt
	// /usr%2F
}
