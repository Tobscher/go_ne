package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/tobscher/kiss/configuration"
)

const (
	pluginPrefix = "kiss-plugin"
)

// StartPlugin starts the plugin with the given name.
// This will try to boot an application called `plugin-<plugin-name>`
//
// This method will return an error when the plugin can not be found
// or the plugin exits with an exit code other than 0.
//
// BUG(tobscher) Find something better than having a loop to process one entry.
func RunTask(t *configuration.Task) error {
	var cmd *exec.Cmd

	for key, plugin := range t.Plugin {
		fmt.Printf("Starting plugin `%v`\n", key)
		command := fmt.Sprintf("%v-%v", pluginPrefix, key)

		cmd = exec.Command(command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Start the plugin
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		err = cmd.Start()
		if err != nil {
			return err
		}

		options, err := json.Marshal(plugin.Options)
		if err != nil {
			return err
		}

		io.WriteString(stdin, string(options))
		io.WriteString(stdin, "\n")

		break
	}

	return cmd.Wait()
}
