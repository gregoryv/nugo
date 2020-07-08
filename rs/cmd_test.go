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
	asJohn := NewAccount("john", 2).Use(NewSystem())
	ls := NewCmd("/bin/ls", "-R", "/")
	ls.Out = os.Stdout
	asJohn.Exec(ls)
	// output:
	// d--xrwxr-xr-x 1 1 bin
	// ----rwxr-xr-x 1 1 bin/ls
	// ----rwxr-xr-x 1 1 bin/mkdir
	// d---rwxr-xr-x 1 1 etc
	// d---rwxr-xr-x 1 1 etc/accounts
	// ----rw-r--r-- 1 1 etc/accounts/anonymous.acc
	// ----rw-r--r-- 1 1 etc/accounts/root.acc
	// drwxrwxrwxrwx 1 1 tmp
}

func ExampleAccount_Exec_lsRecursive() {
	sys := NewSystem()
	// hide directories
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
