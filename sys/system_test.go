package sys

import (
	"fmt"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func ExampleNewSystem() {
	NewSystem().dumprs(os.Stdout)
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /bin
	// d---rwxr-xr-x 1 1 /etc
	// d---rwxr-xr-x 1 1 /etc/accounts
}

func ExampleSystem_Stat() {
	sys := NewSystem()
	_, err := sys.Stat("/etc/accounts", Anonymous)
	fmt.Println(err)
	// output:
	// Stat /etc/accounts uid:0: d---rwxr-xr-x 1 1 exec denied
}

func TestSystem_Stat(t *testing.T) {
	var (
		sys     = NewSystem()
		Stat    = sys.Stat
		ok, bad = asserter.NewMixed(t)
	)
	ok(Stat("/", Anonymous))
	ok(Stat("/bin", Anonymous))
	bad(Stat("/etc", Anonymous))
}
