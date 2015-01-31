package main

import (
	"github.com/tobscher/kiss/plugins/core"
	"github.com/tobscher/kiss/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	command := args.Options["command"]
	*reply = shared.NewResponse(shared.NewCommand("sh -c -e", command))

	return nil
}

func NewShellCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewShellCommand())
	plugin.Serve()
}
