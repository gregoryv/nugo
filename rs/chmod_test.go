package rs

import (
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestChmod(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asAnonymous := Anonymous.Use(sys)
	ok, bad := asserter.NewErrors(t)
	ok(asRoot.Exec("/bin/chmod -m 01755 /tmp"))
	bad(asAnonymous.Exec("/bin/chmod -m 01755 /tmp"))
	bad(asRoot.Exec("/bin/chmod -badflag 01755"))
	bad(asRoot.Exec("/bin/chmod -m 01755"))
	bad(asRoot.Exec("/bin/chmod -m 010000 /tmp"))
}

func ExampleChmod() {
	asRoot := Root.Use(NewSystem())
	asRoot.Fexec(os.Stdout, "/bin/ls", "/tmp") // before
	asRoot.Exec("/bin/chmod -m 01755 /tmp")
	asRoot.Fexec(os.Stdout, "/bin/ls", "/tmp") // after
	// output:
	// drwxrwxrwxrwx 1 1 tmp
	// d--xrwxr-xr-x 1 1 tmp
}

func ExampleChmod_help() {
	asRoot := Root.Use(NewSystem())
	asRoot.Fexec(os.Stdout, "/bin/chmod", "-h")
	// output:
	// Usage of chmod:
	//   -m uint
	//     	mode
}
