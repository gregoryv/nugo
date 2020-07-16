package rs

import (
	"flag"
	"fmt"

	"github.com/gregoryv/nugo"
)

// Ls lists resources
func Ls(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("ls", flag.ContinueOnError)
	longList := flags.Bool("l", false, "use a long listing format")
	recursive := flags.Bool("R", false, "recursive")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		return err
	}
	abspath := flags.Arg(0)
	visitor := func(c *ResInfo, abspath string, w *nugo.Walker) {
		switch {
		case *recursive:
			fmt.Fprintf(cmd.Out, "%s\n", abspath[1:])
		default:
			fmt.Fprintf(cmd.Out, "%s\n", c.Name())
		}
	}
	if *longList {
		visitor = func(c *ResInfo, abspath string, w *nugo.Walker) {
			switch {
			case *recursive:
				fmt.Fprintf(cmd.Out, "%s %s\n", c.node.Seal(), abspath[1:])
			default:
				fmt.Fprintf(cmd.Out, "%s\n", c.node)
			}
		}
	}
	w := NewWalker(cmd.Sys)
	w.SetRecursive(*recursive)
	w.SkipSibling()
	res, err := cmd.Sys.Stat(abspath)
	if err != nil {
		return err
	}
	return w.Walk(res, visitor)
}
