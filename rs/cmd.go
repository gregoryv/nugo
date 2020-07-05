package rs

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gregoryv/nugo"
)

// NewCmd returns a new command.
func NewCmd(abspath string, args ...string) *Cmd {
	return &Cmd{
		Abspath: abspath, Args: args, Out: ioutil.Discard}
}

type Cmd struct {
	Abspath string // of the command
	Args    []string

	// Access to system with a specific account
	Sys *Syscall

	Out io.Writer
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

// ----------------------------------------

type lsCmd struct{}

// Exec
func (me *lsCmd) Exec(c *Cmd) error {
	flags := flag.NewFlagSet("ls", flag.ContinueOnError)
	flags.Usage = func() {}
	flags.SetOutput(ioutil.Discard)
	if err := flags.Parse(c.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	nodes, err := c.Sys.ls(abspath)
	if err != nil {
		return err
	}
	for _, n := range nodes {
		fmt.Fprintf(c.Out, "%s\n", n)
	}
	return nil
}
