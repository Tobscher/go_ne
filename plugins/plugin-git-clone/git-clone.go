package main

import (
	"strings"

	"github.com/tobscher/kiss/plugins/core"
	"github.com/tobscher/kiss/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	directory := args.Options["directory"]

	a := []string{
		"git",
		"clone",
		args.Options["repo"],
		directory,
	}

	cmd := shared.NewCommand(strings.Join(a, " "))
	cmd.Unless(shared.DirectoryExists(directory))

	*reply = shared.NewResponse(cmd)

	return nil
}

func NewCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewCommand())
	plugin.Serve()
}
