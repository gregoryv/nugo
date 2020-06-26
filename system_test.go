package nugo

import "os"

func ExampleSystem_dumprs() {
	sys := NewSystem()
	sys.dumprs(os.Stdout)
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /bin
}
