package sys

import (
	"fmt"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func ExampleSystem_dumprs() {
	sys := NewSystem()
	sys.dumprs(os.Stdout)
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /bin
	// d---rwxr-xr-x 1 1 /etc
	// d---rwxr-xr-x 1 1 /etc/accounts
}

func ExampleSystem_stat() {
	sys := NewSystem()
	_, err := sys.stat("/etc/accounts", Anonymous)
	fmt.Println(err)
	// output:
	// stat /etc/accounts uid:0: d---rwxr-xr-x 1 1 exec denied
}

func TestSystem_stat(t *testing.T) {
	var (
		sys     = NewSystem()
		stat    = sys.stat
		ok, bad = asserter.NewMixed(t)
	)
	ok(stat("/", Anonymous))
	ok(stat("/bin", Anonymous))
	bad(stat("/etc", Anonymous))
}
