package main

import (
	"github.com/tobscher/kiss/plugins/core"
	"github.com/tobscher/kiss/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	*reply = shared.NewResponse(shared.NewCommand("whoami"))

	return nil
}

func NewCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewCommand())
	plugin.Serve()
}
