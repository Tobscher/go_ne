package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tobscher/go_ne/configuration"
)

// NewHostsCommand creates the hosts command.
func NewHostsCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "hosts",
		Short: "Host related option",
		Long:  "Host related option, e.g. list hosts",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	command.AddCommand(listHostsCommand())

	return command
}

func listHostsCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "Print list of available hosts",
		Long:  "Print list of available hosts",
		Run: func(cmd *cobra.Command, args []string) {
			config := configuration.Load(configFile)

			fmt.Println("Hosts:")
			for _, host := range config.Hosts {
				fmt.Printf("  %v\t%v\n", host.Host, host.Description)
			}
		},
	}

	return command
}
