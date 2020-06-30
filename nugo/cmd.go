package nugo

import (
	"flag"
	"io/ioutil"

	"github.com/gregoryv/rs"
)

// NewCmd returns a new command.
func NewCmd(abspath string, args ...string) *Cmd {
	return &Cmd{Abspath: abspath, Args: args}
}

type Cmd struct {
	// todo maybe add a syscall wrapper so commands cannot switch
	// accounts
	Abspath string
	Args    []string

	// Access to system with a specific account
	Sys *Syscall
}

type mkdirCmd struct{}

func (me *mkdirCmd) Exec(c *Cmd) error {
	flags := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.SetOutput(ioutil.Discard)
	mode := flags.Uint("m", 00755, "mode for new directory")
	if err := flags.Parse(c.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	_, err := c.Sys.Mkdir(abspath, rs.NodeMode(*mode))
	return err
}
