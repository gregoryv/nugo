package graph

import "testing"

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
