package rs

import (
	"bytes"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestSyscall_Exec_ls(t *testing.T) {
	var (
		sys     = NewSystem()
		asRoot  = Root.Use(sys)
		ok, bad = asserter.NewErrors(t)
		buf     bytes.Buffer
	)
	bad(asRoot.ExecCmd("/bin/ls", "-xx"))
	bad(asRoot.ExecCmd("/bin/ls", "/nosuch"))

	// ls directory
	ls := NewCmd("/bin/ls", "-R", "/")
	ls.Out = &buf
	ok(asRoot.Exec(ls))
	got := buf.String()
	if got == "" {
		t.Error("missing output")
	}

	// ls file
	asRoot.Save("/tmp/alien", &Alien{Name: "red"})
	buf.Reset()
	ls = NewCmd("/bin/ls", "/tmp/alien")
	ls.Out = &buf
	ok(asRoot.Exec(ls))
	got = buf.String()
	if got == "" {
		t.Error("missing output")
	}
}

func Example_binLs() {
	asRoot := Root.Use(NewSystem())
	ls := NewCmd("/bin/ls", "-R", "/")
	ls.Out = os.Stdout
	asRoot.Exec(ls)
	// output:
	// d--xrwxr-xr-x 1 1 bin
	// ----rwxr-xr-x 1 1 bin/ls
	// ----rwxr-xr-x 1 1 bin/mkdir
	// d---rwxr-xr-x 1 1 etc
	// d---rwxr-xr-x 1 1 etc/accounts
	// drwxrwxrwxrwx 1 1 tmp
}

func ExampleAccount_Exec_lsRecursive() {
	sys := NewSystem()
	// hide etc
	asRoot := Root.Use(sys)
	n, _ := asRoot.stat("/etc")
	n.SetPerm(0)
	// recursive ls
	ls := NewCmd("/bin/ls", "-R", "/")
	ls.Out = os.Stdout
	NewAccount("john", 2).Use(sys).Exec(ls)
	// output:
	// d--xrwxr-xr-x 1 1 bin
	// ----rwxr-xr-x 1 1 bin/ls
	// ----rwxr-xr-x 1 1 bin/mkdir
	// d------------ 1 1 etc
	// drwxrwxrwxrwx 1 1 tmp
}
