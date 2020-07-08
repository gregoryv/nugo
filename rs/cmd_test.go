package rs

import (
	"bytes"
	"os"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestBin_Ls(t *testing.T) {
	asRoot := Root.Use(NewSystem())
	ok, bad := asserter.NewErrors(t)
	bad(asRoot.Exec("/bin/ls", "-xx"))
	bad(asRoot.Exec("/bin/ls", "/nosuch"))
	// ls directory is covered by Examples
	// ls file
	asRoot.Save("/tmp/alien", &Alien{Name: "red"})
	ls := NewCmd("/bin/ls", "/tmp/alien")
	var buf bytes.Buffer
	ls.Out = &buf
	ok(asRoot.ExecCmd(ls))
	got := buf.String()
	if got == "" {
		t.Error("missing output")
	}
}

func ExampleBin_Ls() {
	asJohn := NewAccount("john", 2).Use(NewSystem())
	asJohn.Fexec(os.Stdout, "/bin/ls", "-R", "/")
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

func ExampleAccount_ExecCmd_lsRecursive() {
	sys := NewSystem()
	// hide directories
	asRoot := Root.Use(sys)
	n, _ := asRoot.stat("/etc")
	n.SetPerm(0)
	// recursive ls
	ls := NewCmd("/bin/ls", "-R", "/")
	ls.Out = os.Stdout
	NewAccount("john", 2).Use(sys).ExecCmd(ls)
	// output:
	// d--xrwxr-xr-x 1 1 bin
	// ----rwxr-xr-x 1 1 bin/ls
	// ----rwxr-xr-x 1 1 bin/mkdir
	// d------------ 1 1 etc
	// drwxrwxrwxrwx 1 1 tmp
}
