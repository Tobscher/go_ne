package core

// Command describes a command that should run on
// a remote system
type Command struct {
	name string
	args []string
}

// Name returns the name of the command, e.g. ping
func (c *Command) Name() string {
	return c.name
}

// Args returns the arguments of a command, e.g. -c 3
func (c *Command) Args() []string {
	return c.args
}

// NewCommand creates a new command with args.
func NewCommand(name string, args ...string) (*Command, error) {
	command := Command{
		name: name,
		args: args,
	}

	return &command, nil
}
