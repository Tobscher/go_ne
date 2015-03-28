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
	err := plugin.RunCommand(s.options.Sudo, s.options.Command, s.options.Args...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return 0
}
