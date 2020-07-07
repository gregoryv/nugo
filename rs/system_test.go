package rs

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/fox"
	"github.com/gregoryv/nugo"
)

func TestSystem_SetAuditer(t *testing.T) {
	var (
		buf     bytes.Buffer
		sys     = NewSystem().SetAuditer(fox.NewSyncLog(&buf))
		asRoot  = Root.Use(sys) // use after the auditer is set
		asJohn  = NewAccount("john", 2).Use(sys)
		ok, bad = asserter.NewErrors(t)
	)
	bad(asJohn.ExecCmd("/bin/mkdir", "/etc/s"))
	ok(asJohn.ExecCmd("/bin/mkdir", "/tmp/s"))
	ok(asRoot.ExecCmd("/bin/mkdir", "/etc/x"))
	if buf.String() == "" {
		t.Error("expected audit")
	}
	if !strings.Contains(buf.String(), "ERR") {
		t.Error(buf.String())
	}
}

func TestSystem_rootNode(t *testing.T) {
	var (
		sys = NewSystem()
		_   = sys.mount(nugo.NewRoot("/mnt"))
		_   = sys.mount(nugo.NewRoot("/mnt/usb"))
	)
	if rn := sys.rootNode("/mnt/usb/some/path"); rn.Name() != "/mnt/usb" {
		t.Fail()
	}
	if rn := sys.rootNode("/nosuch/dir"); rn.Name() != "/" {
		t.Fail()
	}
}

func TestSystem_mount(t *testing.T) {
	var (
		sys     = NewSystem()
		ok, bad = asserter.NewErrors(t)
	)
	bad(sys.mount(nugo.NewRoot("/")))
	ok(sys.mount(nugo.NewRoot("/mnt/usb")))
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
