package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tobscher/go_ne/configuration"
)

// NewTasksCommand creates the tasks command.
func NewTasksCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "tasks",
		Short: "Task related option",
		Long:  "Task related option, e.g. list tasks",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	command.AddCommand(listTasksCommand())

	return command
}

func listTasksCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "Print list of available tasks",
		Long:  "Print list of available tasks",
		Run: func(cmd *cobra.Command, args []string) {
			config := configuration.Load(configFile)

			fmt.Println("Global tasks:")
			for _, task := range config.Tasks {
				fmt.Printf("  %v\t%v\n", task.Task, task.Description)
			}

			for _, host := range config.Hosts {
				if len(host.Tasks) == 0 {
					continue
				}

				fmt.Printf("\n%v tasks:\n", host.Host)

				for _, task := range host.Tasks {
					fmt.Printf("  %v\t%v\n", task.Task, task.Description)
				}
			}
		},
	}

	return command
}
