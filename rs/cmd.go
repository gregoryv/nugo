package rs

import (
	"flag"
	"fmt"

	"github.com/gregoryv/nugo"
)

// Mkacc
func Mkacc(cmd *Cmd) ExecErr {
	flags := flag.NewFlagSet("mkacc", flag.ContinueOnError)
	uid := flags.Int("uid", -1, "uid of the new account")
	gid := flags.Int("gid", -1, "gid of the new account")
	flags.SetOutput(cmd.Out)
	if err := flags.Parse(cmd.Args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}
	name := flags.Arg(0)
	if name == "" {
		return fmt.Errorf("missing account name")
	}
	if *uid < 2 {
		return fmt.Errorf("invalid uid")
	}
	if *gid < 2 {
		return fmt.Errorf("invalid gid")
	}
	acc := NewAccount(name, *uid)
	acc.groups[0] = *gid
	return cmd.Sys.AddAccount(acc)
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
