package graph

import "os"

func ExampleSystem_ls() {
	sys := NewSystem()
	sys.dumprs(os.Stdout)
	// output:
	// /
	// /bin
}
