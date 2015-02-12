package main

import (
	"strings"

	"github.com/tobscher/kiss/plugins/shared"
)

type Command struct{}

func NewCommand() *Command {
	return new(Command)
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	var commands []shared.Command

	update := shared.ExtractTruthy(args.Options["update"])
	if update {
		commands = append(commands, updateCommand())
	}

	commands = append(commands, installCommand(args))

	*reply = shared.NewResponse(commands...)

	return nil
}

func updateCommand() shared.Command {
	cmd := []string{
		"sudo", // only if sudo true
		"apt-get",
		"update",
		"-y",
	}

	return shared.NewCommand(strings.Join(cmd, " "))
}

func installCommand(args shared.Args) shared.Command {
	packages := shared.ExtractOptions(args.Options["packages"])

	cmd := []string{
		"sudo", // only if sudo true
		"apt-get",
		"install",
		"-y",
		strings.Join(packages, " "),
	}

	return shared.NewCommand(strings.Join(cmd, " "))
}
