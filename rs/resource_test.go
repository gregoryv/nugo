package rs

import (
	"bytes"
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/nugo"
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

func TestResource_SetSource(t *testing.T) {
	var (
		rw       = newResource(nugo.NewNode("x"), OpRead|OpWrite)
		readOnly = newResource(nugo.NewNode("x"), OpRead)
		ok, bad  = asserter.NewErrors(t)
	)
	ok(rw.SetSource("string"))
	ok(rw.SetSource([]byte("bytes")))
	bad(rw.SetSource(1))
	bad(readOnly.SetSource("read only"))
}

func TestResource_Read(t *testing.T) {
	var (
		ok, bad = asserter.NewMixed(t)
		b       = make([]byte, 10)
	)
	r := &Resource{buf: bytes.NewBufferString("hello")}
	ok(r.Read(b))
	r = &Resource{}
	bad(r.Read(b))
}
