package rs

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
)

func TestResInfo_Name(t *testing.T) {
	var (
		rif, _ = NewSystem().Use(Root).Stat("/")
	)
	if rif.Name() != "/" {
		t.Error("name failed")
	}
}

func TestResInfo_IsDir(t *testing.T) {
	var (
		asRoot  = NewSystem().Use(Root)
		dir, _  = asRoot.Stat("/")
		file, _ = asRoot.Stat("/bin/mkdir")
		ok, bad = asserter.NewErrors(t)
	)
	ok(dir.IsDir())
	bad(file.IsDir())
}

func TestResource_Read(t *testing.T) {
	var (
		ok, bad = asserter.NewMixed(t)
		b       = make([]byte, 10)
	)
	r := &Resource{op: OpRead, buf: bytes.NewBufferString("hello")}
	ok(r.Read(b))
	r = &Resource{op: OpRead}
	bad(r.Read(b))
}
