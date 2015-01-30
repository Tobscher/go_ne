package shared

import "github.com/tobscher/go_ne/configuration"

// Args describes the args for the remote call.
type Args struct {
	Environment []string
	Args        []string
	Options     configuration.OptionCollection
}

// Command describes the command and its args to run.
type Command struct {
	Name string
	Args []string
}

// Response describes the response for the remote call.
type Response struct {
	Commands []Command
}

// Responder is the interface which defines a plugin.
type Responder interface {
	Execute(args Args, reply *Response) error
}

// NewCommand creates a new Command.
func NewCommand(name string, args ...string) Command {
	return Command{
		Name: name,
		Args: args,
	}
}

// NewResponse creates a new response with the given commands.
func NewResponse(commands ...Command) Response {
	return Response{Commands: commands}
}
