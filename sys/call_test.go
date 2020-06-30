package sys

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestSyscall_Mkdir(t *testing.T) {
	var (
		asRoot      = Syscall{acc: Root, sys: NewSystem()}
		asAnonymous = Syscall{acc: Anonymous, sys: NewSystem()}
		ok, bad     = asserter.NewMixed(t)
	)
	ok(asRoot.Mkdir("/adir", 0))

	// parent directory missing
	bad(asRoot.Mkdir("/nosuch/whatever", 0))

	// permission inadequate
	bad(asAnonymous.Mkdir("/whatever", 0))
}
