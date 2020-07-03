package rs

import (
	"fmt"
	"testing"

	"github.com/gregoryv/asserter"
)

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

func TestSyscall_Open_Source(t *testing.T) {
	var (
		asRoot  = Root.Use(NewSystem())
		_, _    = asRoot.Install("/string", "a string", 00644)
		_, _    = asRoot.Install("/byteSlice", []byte("bytes"), 00644)
		_, _    = asRoot.Install("/int", 1, 00644)
		ok, bad = asserter.NewMixed(t)
	)
	ok(asRoot.Open("/string"))
	ok(asRoot.Open("/byteSlice"))
	ok(asRoot.Open("/bin/mkdir"))
	bad(asRoot.Open("/int"))
}

func TestSyscall_Open(t *testing.T) {
	var (
		sys         = NewSystem()
		asRoot      = Root.Use(sys)
		_, _        = asRoot.Install("/existing.dat", "content", 00644)
		asAnonymous = Anonymous.Use(sys)
		ok, bad     = asserter.NewMixed(t)
	)
	// owner has read permission on newly created resources
	ok(asRoot.Open("/existing.dat"))
	// missing resource
	bad(asRoot.Open("/nosuch"))
	// inadequate permission
	bad(asAnonymous.Open("/existing.dat"))
	// write to read only
	res, _ := asRoot.Open("/existing.dat")
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

func TestSyscall_Exec(t *testing.T) {
	var (
		Exec    = Root.Use(NewSystem()).Exec
		ok, bad = asserter.NewErrors(t)
	)
	ok(Exec(NewCmd("/bin/mkdir", "/tmp")))
	// Node not found
	bad(Exec(NewCmd("/bin/nosuch/mkdir", "/tmp")))
	// Resource not type Executable
	bad(Exec(NewCmd("/bin")))
	// Bad flag
	bad(Exec(NewCmd("/bin/mkdir", "-nosuch")))
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
