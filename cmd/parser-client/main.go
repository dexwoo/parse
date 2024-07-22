package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"parser/gen/pb"
)

func main() {
	cmd := &command{
		stdout: os.Stdout,
	}
	if err := cmd.run(); err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err)
	}
}

type command struct {
	stdout io.Writer
}

var u = flag.String("url", "localhost:8888", "parser server URL")

func (cmd *command) run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if len(os.Args) < 2 {
		return cmdError("command required")
	}

	switch subcmd := os.Args[1]; subcmd {
	case "get-block":
		f := newFlagSet(subcmd)
		if err := f.Parse(os.Args[2:]); err != nil {
			return flagError(f, err.Error(), subcmd)
		}
		if u == nil || *u == "" {
			return flagError(f, "-url flag required", subcmd)
		}
		c := pb.NewParserClient(*u)
		return c.GetCurrentBlock(ctx)
	default:
		return cmdError(fmt.Sprintf("unknown command %s", subcmd))
	}
}

func newFlagSet(name string) *flag.FlagSet {
	f := flag.NewFlagSet(name, flag.ContinueOnError)
	f.SetOutput(io.Discard)
	return f
}

func cmdError(msg string) error {
	return fmt.Errorf(`%s

Usage: parser <command>

Available Commands:
  get-block    get current block
`, msg)
}

func flagError(f *flag.FlagSet, msg, subcmd string) error {
	var b bytes.Buffer
	f.SetOutput(&b)
	f.PrintDefaults()
	return fmt.Errorf(`%s

Usage: parser %s [flags]

Available Flags:
%s`, msg, subcmd, b.String())
}
