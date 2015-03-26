package main

import (
	"log"

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
		log.Fatal(err)
	}

	return 0
}
