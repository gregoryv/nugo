package rs

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gregoryv/nugo"
)

var (
	Anonymous = NewAccount("anonymous", 0)
	Root      = NewAccount("root", 1)
)

// NewAccount returns a new account with the given uid as both uid and
// group id.
func NewAccount(name string, uid int) *Account {
	return &Account{
		name:   name,
		uid:    uid,
		groups: []int{uid},
	}
}

type Account struct {
	name string
	uid  int

	mu     sync.Mutex
	groups []int
}

// gid returns the first group id of the account
func (my *Account) gid() int { return my.groups[0] }

// todo hide as command
func (me *Account) AddGroup(gid int) {
	for _, id := range me.groups {
		if id == gid {
			return
		}
	}
	me.mu.Lock()
	me.groups = append(me.groups, gid)
	me.mu.Unlock()
}

// todo hide as command
func (me *Account) DelGroup(gid int) {
	for i, id := range me.groups {
		if id == gid {
			me.mu.Lock()
			me.groups = append(me.groups[:i], me.groups[i+1:]...)
			me.mu.Unlock()
			return
		}
	}
}

// Use returns a Syscall struct for accessing the system.
func (me *Account) Use(sys *System) *Syscall {
	return &Syscall{System: sys, acc: me}
}

func (me *Account) owns(id int) bool { return me.uid == id }

// permitted returns error if account does not have operation
// permission to the given seal.
func (my *Account) permitted(op operation, seal *nugo.Seal) error {
	if my.uid == Root.uid {
		return nil
	}
	n, u, g, o := op.Modes()
	switch {
	case my.uid == 0 && (seal.Mode&n == n): // anonymous
	case my.uid == seal.UID && (seal.Mode&u == u): // owner
	case my.member(seal.GID) && (seal.Mode&g == g): // group
	case my.uid > 0 && seal.Mode&o == o: // other
	default:
		return fmt.Errorf("%v %v denied", seal, op)
	}
	return nil
}

var ErrPermissionDenied = errors.New("permission denied")

func (my *Account) member(gid int) bool {
	for _, id := range my.groups {
		if id == gid {
			return true
		}
	}
	return false
}
