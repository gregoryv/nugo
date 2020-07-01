package rs

import (
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
