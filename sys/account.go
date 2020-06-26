package sys

import "sync"

// NewAccount returns a new account with the given uid as both uid and
// group id.
func NewAccount(username string, uid int) *Account {
	return &Account{
		username: username,
		uid:      uid,
		groups:   []int{uid},
	}
}

type Account struct {
	username string
	uid      int

	mu     sync.Mutex
	groups []int
}

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

func (me *Account) owns(id int) bool { return me.uid == id }
