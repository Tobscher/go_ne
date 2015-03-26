package core

import "github.com/tobscher/kiss/configuration"

// RunTask uses the given runner to execute a task.
func RunTask(runner Runner, t *configuration.Task) error {
	defer runner.Close()

	runner.Run(t)

	return nil
}
