package sys

import (
	"fmt"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestSyscall_Mkdir(t *testing.T) {
	var (
		asRoot      = Syscall{System: NewSystem(), acc: Root}
		asAnonymous = Syscall{System: NewSystem(), acc: Anonymous}
		ok, bad     = asserter.NewMixed(t)
	)
	ok(asRoot.Mkdir("/adir", 0))

	// parent directory missing
	bad(asRoot.Mkdir("/nosuch/whatever", 0))

	// permission inadequate
	bad(asAnonymous.Mkdir("/whatever", 0))
}

func TestSyscall_Exec(t *testing.T) {
	var (
		sys     = &Syscall{System: NewSystem(), acc: Root}
		Exec    = sys.Exec
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
	sys := &Syscall{System: NewSystem(), acc: Anonymous}
	_, err := sys.Stat("/etc/accounts")
	fmt.Println(err)
	// output:
	// Stat /etc/accounts uid:0: d---rwxr-xr-x 1 1 exec denied
}

func TestSystem_Stat(t *testing.T) {
	var (
		sys     = &Syscall{System: NewSystem(), acc: Anonymous}
		Stat    = sys.Stat
		ok, bad = asserter.NewMixed(t)
	)
	ok(Stat("/"))
	ok(Stat("/bin"))
	bad(Stat("/etc"))
	bad(Stat("/nothing"))
}

func TestSystem_Install(t *testing.T) {
	var (
		asAnonymous = &Syscall{System: NewSystem(), acc: Anonymous}
		asRoot      = &Syscall{System: NewSystem(), acc: Root}
		ok, bad     = asserter.NewMixed(t)
	)
	ok(asRoot.Install("/bin/x", nil, 0))
	bad(asRoot.Install("/bin/nosuchdir/x", nil, 0))
	bad(asAnonymous.Install("/bin/x", nil, 0))
}
