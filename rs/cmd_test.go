package rs

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/asserter"
)

var (
	sys    = NewSystem()
	asRoot = Root.Use(sys)
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
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asRoot.Exec("/bin/chmod -m 0 /tmp")
	asRoot.Fexec(os.Stdout, "/bin/ls", "/tmp")
	// output:
	// d------------ 1 1 tmp
}

func TestLs(t *testing.T) {
	sys := NewSystem()
	asRoot := Root.Use(sys)
	asJohn := NewAccount("john", 2).Use(sys)
	var buf bytes.Buffer
	ok, bad := asserter.NewErrors(t)
	bad(asRoot.Exec("/bin/ls -xx"))
	bad(asRoot.Exec("/bin/ls /nosuch"))
	// ls directory is covered by Examples
	// ls file
	asRoot.Fexec(&buf, "/bin/ls", "/etc/accounts/root.acc")
	exp := "----rw-r--r-- 1 1 root.acc\n"
	assert := asserter.New(t)
	assert().Equals(buf.String(), exp)

	// only list accessible
	buf.Reset()
	n, _ := asRoot.stat("/etc")
	n.SetPerm(0)
	ok(asJohn.Fexec(&buf, "/bin/ls", "-R", "/"))
	if strings.Contains(buf.String(), "/etc/accounts") {
		t.Error("listed /etc")
	}
}

func ExampleLs() {
	Anonymous.Use(NewSystem()).Fexec(os.Stdout, "/bin/ls", "/")
	// output:
	// d--xrwxr-xr-x 1 1 /
	// d--xrwxr-xr-x 1 1 bin
	// d---rwxr-xr-x 1 1 etc
	// drwxrwxrwxrwx 1 1 tmp
}
