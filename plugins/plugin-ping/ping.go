package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tobscher/kiss/plugins/core"
	"github.com/tobscher/kiss/plugins/shared"
)

type Command struct {
}

func (t *Command) Execute(args shared.Args, reply *shared.Response) error {
	cmd := []string{
		"ping",
	}

	if count, ok := args.Options["count"]; ok {
		c, _ := strconv.ParseFloat(count, 64)
		cmd = append(cmd, fmt.Sprintf("-c %v", c))
	}

	cmd = append(cmd, shared.ExtractString(args.Options["url"]))
	*reply = shared.NewResponse(shared.NewCommand(strings.Join(cmd, " ")))

	return nil
}

func NewEnvCommand() *Command {
	return new(Command)
}

func main() {
	plugin.Register(NewEnvCommand())
	plugin.Serve()
}
