package graph

import "os"

func ExampleSystem_ls() {
	sys := NewSystem()
	sys.ls(os.Stdout)
	// output:
	// /
	// /bin
}
