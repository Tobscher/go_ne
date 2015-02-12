package main

import "github.com/tobscher/kiss/plugins/shared"

type Command struct{}

func NewCommand() *Command {
	return new(Command)
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	command := args.Options["command"]
	*reply = shared.NewResponse(shared.NewCommand("sh -c -e", command))

	return nil
}
