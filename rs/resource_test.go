package rs

import (
	"bytes"
	"io"
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
		rw       = &Resource{node: nugo.NewNode("x")}
		readOnly = &Resource{readOnly: true, node: nugo.NewNode("x")}
		ok, bad  = asserter.NewErrors(t)
	)
	ok(rw.SetSource(1))
	bad(readOnly.SetSource(2))
}

func TestResource_Read(t *testing.T) {
	var (
		n      = nugo.NewNode("node")
		r      = &Resource{readOnly: true, node: n}
		assert = asserter.New(t)
		ok, _  = asserter.NewMixed(t)
		buf    bytes.Buffer
	)
	// Read byte slice
	n.SetSource([]byte("hello"))
	ok(io.Copy(&buf, r))
	assert().Equals(buf.String(), "hello")
	r.Close()
	// Read string
	n.SetSource("world")
	buf.Reset()
	ok(io.Copy(&buf, r))
	assert().Equals(buf.String(), "world")
	r.Close()
}
