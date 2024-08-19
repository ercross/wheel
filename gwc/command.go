package main

import (
	"io"
	"strings"
)

type command struct {
	programName string
	options     []commandOptions
}

type commandOptions struct {
}

func parseArgs(args []string) (command, error) {
	return command{programName: args[0]}, nil
}

func extractInput(input string) (io.Reader, error) {
	return strings.NewReader(input), nil
}
