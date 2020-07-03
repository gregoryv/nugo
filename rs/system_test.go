package rs

import (
	"os"
)

func ExampleNewSystem() {
	NewSystem().dumprs(os.Stdout)
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /bin
	// ----rwxr-xr-x 1 1 /bin/mkdir
	// d---rwxr-xr-x 1 1 /etc
	// d---rwxr-xr-x 1 1 /etc/accounts
	// drwxrwxrwxrwx 1 1 /tmp
}
