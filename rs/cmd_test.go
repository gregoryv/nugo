package rs

import (
	"bytes"
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
	ls := NewCmd("/bin/ls", "/")
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
