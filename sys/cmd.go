package sys

// NewCmd returns a new command.
func NewCmd(abspath string, args ...string) *Cmd {
	return &Cmd{abspath: abspath, args: args}
}

type Cmd struct {
	// todo maybe add a syscall wrapper so commands cannot switch
	// accounts
	abspath string
	args    []string
}

// Run the command on the given system using a specific account.
func (me *Cmd) Run(sys *System, acc *Account) error {
	return sys.Exec(me, acc)
}
