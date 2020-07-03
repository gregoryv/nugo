package rs

import (
	"fmt"
	"os"
)

func Example_newSystem() {
	NewSystem().dumprs(os.Stdout)
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 /bin
	// ----rwxr-xr-x 1 1 /bin/mkdir
	// d---rwxr-xr-x 1 1 /etc
	// d---rwxr-xr-x 1 1 /etc/accounts
	// drwxrwxrwxrwx 1 1 /tmp
}

func Example_saveAndLoadResource() {
	var (
		sys    = NewSystem()
		asRoot = Root.Use(sys)
	)
	asRoot.ExecCmd("/bin/mkdir", "/tmp/aliens")
	asRoot.Save("/tmp/aliens/green.gob", &Alien{Name: "Mr Green"})

	var alien Alien
	asRoot.Load(&alien, "/tmp/aliens/green.gob")
	fmt.Printf("%#v", alien)
	// output:
	// rs.Alien{Name:"Mr Green"}
}
