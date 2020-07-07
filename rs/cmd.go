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
func (me *lsCmd) Exec(cmd *Cmd) error {
	flags := flag.NewFlagSet("ls", flag.ContinueOnError)
	recursive := flags.Bool("R", false, "recursive")
	flags.Usage = func() {}
	flags.SetOutput(ioutil.Discard)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	visitor := func(p, c *ResInfo, abspath string, w *nugo.Walker) {
		switch {
		case p == nil:
			// only when ls -al
			// fmt.Fprintf(cmd.Out, "%s .\n", c.node.Seal())
		case *recursive:
			fmt.Fprintf(cmd.Out, "%s %s\n", c.node.Seal(), abspath[1:])
		default:
			fmt.Fprintf(cmd.Out, "%s\n", c.node)
		}
	}
	return cmd.Sys.Walk(abspath, *recursive, visitor)
}
