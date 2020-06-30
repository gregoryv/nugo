package rs

import "testing"

func TestResInfo_Name(t *testing.T) {
	var (
		rif, _ = NewSystem().Use(Root).Stat("/")
	)
	if rif.Name() != "/" {
		t.Error("name failed")
	}
}
