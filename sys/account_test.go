package sys

import (
	"testing"

	"github.com/gregoryv/asserter"
	"github.com/gregoryv/rs"
)

func TestAccount_AddGroup(t *testing.T) {
	acc := NewAccount("root", 1)
	acc.AddGroup(2)
	acc.AddGroup(2) // nop, already there
	if len(acc.groups) != 2 {
		t.Fail()
	}
}

func TestAccount_DelGroup(t *testing.T) {
	acc := NewAccount("root", 1)
	acc.AddGroup(2)
	acc.DelGroup(2)
	if len(acc.groups) != 1 {
		t.Fail()
	}
}

func TestAccount_owns(t *testing.T) {
	acc := NewAccount("root", 1)
	if acc.owns(2) {
		t.Error("uid 1 owns uid 2")
	}
}

func TestAccount_PermittedAnonymous(t *testing.T) {
	var (
		ok, bad = asserter.NewErrors(t)
		perm    = Anonymous.Permitted
	)
	ok(perm(OpRead, &rs.Seal{1, 1, 07000}))
	ok(perm(OpRead, &rs.Seal{1, 1, 04000}))
	ok(perm(OpWrite, &rs.Seal{1, 1, 02000}))
	ok(perm(OpExec, &rs.Seal{1, 1, 01000}))
	bad(perm(OpExec, &rs.Seal{1, 1, 02000}))
	bad(perm(OpExec, &rs.Seal{1, 1, 00000}))
}

func TestAccount_PermittedRoot(t *testing.T) {
	var (
		ok, _ = asserter.NewErrors(t)
		perm  = Root.Permitted
	)
	// root is special in that it always has full access
	ok(perm(OpRead, &rs.Seal{1, 1, 00000}))
	ok(perm(OpWrite, &rs.Seal{1, 2, 00000}))
	ok(perm(OpExec, &rs.Seal{0, 0, 00000}))
}

func TestAccount_PermittedOther(t *testing.T) {
	var (
		ok, _ = asserter.NewErrors(t)
		perm  = NewAccount("john", 2).Permitted
	)
	// root is special in that it always has full access
	ok(perm(OpRead, &rs.Seal{2, 2, 00400}))
	ok(perm(OpRead, &rs.Seal{3, 2, 00040}))
	ok(perm(OpRead, &rs.Seal{1, 1, 00004}))
}
