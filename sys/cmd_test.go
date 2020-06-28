package sys

import (
	"testing"

	"github.com/gregoryv/asserter"
)

func TestCmd_Run(t *testing.T) {
	var (
		sys     = NewSystem()
		ok, bad = asserter.NewErrors(t)
	)
	ok(NewCmd("/bin/mkdir", "/tmp").Run(sys, Root))

	// Node not found
	bad(NewCmd("/bin/nosuch/mkdir", "/tmp").Run(sys, Root))

	// Resource not type Executable
	bad(NewCmd("/bin").Run(sys, Root))

	// Bad flag
	bad(NewCmd("/bin/mkdir", "-nosuch").Run(sys, Root))
}
