package main

import (
	"fmt"
	"os"

	"github.com/tobscher/kiss/plugins/core"
)

type Shell struct {
	options Options
}

func NewShell(options Options) *Shell {
	shell := Shell{
		options: options,
	}

	return &shell
}

func (s *Shell) Run() int {
	err := plugin.RunCommand(s.options.Command, s.options.Sudo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return 0
}
