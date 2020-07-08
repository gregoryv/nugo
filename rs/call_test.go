package rs

import (
	"fmt"
	"testing"

	"github.com/gregoryv/asserter"
)

// test struct
type Alien struct {
	Name string
}

func TestSyscall_RemoveAll(t *testing.T) {
	var (
		sys         = NewSystem()
		asRoot      = Root.Use(sys)
		asAnonymous = Anonymous.Use(sys)
		_           = asRoot.SaveAs("/tmp/alien", &Alien{Name: "RemoveAll"})
		ok, bad     = asserter.NewErrors(t)
	)
	ok(asRoot.RemoveAll("/tmp/alien"))
	bad(asRoot.RemoveAll("/tmp/nosuch"))
	bad(asAnonymous.RemoveAll("/etc/accounts"))
}

func TestSyscall_Load(t *testing.T) {
	var (
		asRoot  = Root.Use(NewSystem())
		alien   = Alien{Name: "Mr green"}
		got     Alien
		ok, bad = asserter.NewErrors(t)
		assert  = asserter.New(t)
	)
	ok(asRoot.Save("/thing.gob", &alien))
	ok(asRoot.Load(&got, "/thing.gob"))
	assert(got.Name == "Mr green").Errorf("%v", got)
	bad(asRoot.Load(&got, "/nosuch"))
	bad(asRoot.Load(&got, "/bin/mkdir"))
}

func TestSyscall_Save(t *testing.T) {
	var (
		asRoot  = Root.Use(NewSystem())
		thing   = struct{ Name string }{"thing"}
		ok, bad = asserter.NewErrors(t)
	)
	ok(asRoot.Save("/thing.gob", &thing))
	ok(asRoot.Save("/thing.gob", &thing)).Error("First fix Create")
	bad(asRoot.Save("/nosuch/thing.gob", &thing))
	bad(asRoot.Save("/", thing))
}

func TestSyscall_SaveAs(t *testing.T) {
	var (
		asRoot  = Root.Use(NewSystem())
		thing   = struct{ Name string }{"thing"}
		ok, bad = asserter.NewErrors(t)
	)
	ok(asRoot.SaveAs("/thing.gob", &thing))
	bad(asRoot.SaveAs("/thing.gob", &thing))
	bad(asRoot.SaveAs("/nosuch/thing.gob", &thing))
	bad(asRoot.SaveAs("/", thing))
}

func TestSyscall_Open(t *testing.T) {
	var (
		sys         = NewSystem()
		asRoot      = Root.Use(sys)
		asAnonymous = Anonymous.Use(sys)
		_           = asRoot.Save("/tmp/alien.gob", &Alien{Name: "x"})
		ok, bad     = asserter.NewMixed(t)
	)
	// owner has read permission on newly created resources
	ok(asRoot.Open("/tmp/alien.gob"))
	// missing resource
	bad(asRoot.Open("/nosuch"))
	// inadequate permission
	bad(asAnonymous.Open("/tmp/alien.gob"))
	// write to read only
	res, _ := asRoot.Open("/tmp/alien.gob")
	bad(res.Write([]byte("")))
}

func TestSyscall_Create(t *testing.T) {
	var (
		sys         = NewSystem()
		asRoot      = Root.Use(sys)
		asAnonymous = Anonymous.Use(sys)
		ok, bad     = asserter.NewMixed(t)
	)
	// new resource
	res, err := asRoot.Create("/file")
	ok(res, err)
	// write over existing
	ok(asRoot.Create("/file"))
	// write only
	bad(res.Read([]byte{}))
	bad(asRoot.Create("/"))
	bad(asAnonymous.Create("/file"))
}

func TestSyscall_Mkdir(t *testing.T) {
	var (
		asRoot      = Root.Use(NewSystem())
		asAnonymous = Anonymous.Use(NewSystem())
		ok, bad     = asserter.NewMixed(t)
	)
	ok(asRoot.Mkdir("/adir", 0))
	// parent directory missing
	bad(asRoot.Mkdir("/nosuch/whatever", 0))
	// inadequate permission
	bad(asAnonymous.Mkdir("/whatever", 0))
}

func TestSyscall_ExecCmd(t *testing.T) {
	var (
		ExecCmd = Root.Use(NewSystem()).ExecCmd
		ok, bad = asserter.NewErrors(t)
	)
	ok(ExecCmd(NewCmd("/bin/mkdir", "/tmp")))
	// Node not found
	bad(ExecCmd(NewCmd("/bin/nosuch/mkdir", "/tmp")))
	// Resource not type Executable
	bad(ExecCmd(NewCmd("/bin")))
	// Bad flag
	bad(ExecCmd(NewCmd("/bin/mkdir", "-nosuch")))
}

func ExampleSyscall_Stat() {
	sys := Anonymous.Use(NewSystem())
	_, err := sys.Stat("/etc/accounts")
	fmt.Println(err)
	// output:
	// Stat /etc/accounts uid:0: d---rwxr-xr-x 1 1 exec denied
}

func TestSystem_Stat(t *testing.T) {
	var (
		Stat    = Anonymous.Use(NewSystem()).Stat
		ok, bad = asserter.NewMixed(t)
	)
	ok(Stat("/"))
	ok(Stat("/bin"))
	bad(Stat("/etc"))
	bad(Stat("/nothing"))
}

func TestSystem_Install(t *testing.T) {
	var (
		sys         = NewSystem()
		asRoot      = Root.Use(sys)
		asAnonymous = Anonymous.Use(sys)
		ok, bad     = asserter.NewMixed(t)
	)
	ok(asRoot.Install("/bin/x", nil, 0))
	bad(asRoot.Install("/bin/nosuchdir/x", nil, 0))
	bad(asAnonymous.Install("/bin/x", nil, 0))
}
