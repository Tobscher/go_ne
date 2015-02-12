package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tobscher/kiss/configuration"
)

// NewRolesCommand creates the roles command.
func NewRolesCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "roles",
		Short: "Role related option",
		Long:  "Role related option, e.g. list roles",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	command.AddCommand(listRolesCommand())

	return command
}

func listRolesCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "Print list of available roles",
		Long:  "Print list of available roles",
		Run: func(cmd *cobra.Command, args []string) {
			config := configuration.Load(configFile)

			fmt.Println("Roles:")
			for name, role := range config.Roles {
				extends := ""
				if len(role.With) > 0 {
					extends = fmt.Sprintf(" (includes %v)", strings.Join(role.With, ","))
				}

				fmt.Printf("  %v\t%v%v\n", name, role.Description, extends)
			}
		},
	}

	return command
}
