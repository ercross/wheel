package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if args[0] != programName {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Errorf("unknown command `%s`", os.Args[0]))
	}

	r, err := run(args[1:])
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
	_, _ = fmt.Fprintln(os.Stdout, r)
}

func run(args []string) (r result, err error) {
	cmd, err := parseArgs(args)
	if err != nil {
		return r, err
	}

	return cmd.process()
}
