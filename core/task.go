package core

import (
	"os"

	"github.com/tobscher/go_ne/configuration"
	"github.com/tobscher/go_ne/plugins/shared"
)

// Task is the interface which describes tasks
// that should be run on the remote system.
type Task interface {
	Name() string
	Args() []string
}

// RunTask uses the given runner to execute a task.
func RunTask(runner Runner, t *configuration.Task) error {
	defer StopAllPlugins()
	defer runner.Close()
	var commands []*Command

	for name, plugin := range t.Plugin {
		// Load plugin
		p, err := GetPlugin(name)
		if err != nil {
			return err
		}

		pluginArgs := shared.Args{
			Environment: os.Environ(),
			Options:     plugin.Options,
		}

		commands, err = p.GetCommands(pluginArgs)
		if err != nil {
			return err
		}

		for _, c := range commands {
			err = runner.Run(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
