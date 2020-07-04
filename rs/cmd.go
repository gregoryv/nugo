package rs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gregoryv/nugo"
)

// NewCmd returns a new command.
func NewCmd(abspath string, args ...string) *Cmd {
	return &Cmd{Abspath: abspath, Args: args}
}

type Cmd struct {
	Abspath string // of the command
	Args    []string

	// Access to system with a specific account
	Sys *Syscall
}

// String
func (me *Cmd) String() string {
	return fmt.Sprintf("%s %s", me.Abspath, strings.Join(me.Args, " "))
}

// ----------------------------------------

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
	_, err := c.Sys.mkdir(abspath, nugo.NodeMode(*mode))
	return err
}
