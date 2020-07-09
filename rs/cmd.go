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

// String returns the command with its arguments
func (me *Cmd) String() string {
	return fmt.Sprintf("%s %s", me.Abspath, strings.Join(me.Args, " "))
}

// ----------------------------------------

type ExecFunc func(*Cmd) ExecErr

func (me ExecFunc) Exec(cmd *Cmd) ExecErr { return me(cmd) }

type ExecErr error

type Executable interface {
	Exec(*Cmd) ExecErr
}

// Chmod
func Chmod(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("chmod", flag.ContinueOnError)
	mode := flags.Uint("m", 0, "mode")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	if abspath == "" {
		return fmt.Errorf("missing abspath")
	}
	return cmd.Sys.SetMode(abspath, Mode(*mode))
}

// Mkdir creates directories
func Mkdir(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("mkdir", flag.ContinueOnError)
	flags.SetOutput(cmd.Out)
	mode := flags.Uint("m", 00755, "mode for new directory")
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	_, err := cmd.Sys.Mkdir(abspath, Mode(*mode))
	return err
}

// Ls lists resources
func Ls(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("ls", flag.ContinueOnError)
	recursive := flags.Bool("R", false, "recursive")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	visitor := func(p, c *ResInfo, abspath string, w *nugo.Walker) {
		switch {
		case p == nil:
			fmt.Fprintf(cmd.Out, "%s %s\n", c.node.Seal(), c.Name())
		case *recursive:
			fmt.Fprintf(cmd.Out, "%s %s\n", c.node.Seal(), abspath[1:])
		default:
			fmt.Fprintf(cmd.Out, "%s\n", c.node)
		}
	}
	return cmd.Sys.Walk(abspath, *recursive, visitor)
}
