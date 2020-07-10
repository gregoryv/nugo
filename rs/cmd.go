package rs

import (
	"flag"
	"fmt"

	"github.com/gregoryv/nugo"
)

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
