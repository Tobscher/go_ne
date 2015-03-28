package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tobscher/kiss/configuration"
	"github.com/tobscher/kiss/core"
	"github.com/tobscher/kiss/logging"
)

// NewRunRoleCommand creates a new command to run a task.
func NewRunRoleCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "run-role",
		Short: "Run a role",
		Long:  "Run a role on all remote systems with that role.",
		Run:   runRoleRun,
	}
	command.Flags().StringVar(&configFile, "config", ".kiss.yml", "path to config file")
	command.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	return command
}

func runRoleRun(cmd *cobra.Command, args []string) {
	if verbose {
		logger.SetLevel(logging.DEBUG)
		core.SetLogLevel(logging.DEBUG)
	}

	config := configuration.Load(configFile)

	if len(args) < 1 {
		fmt.Printf("Error: Expected role-name: kiss run-role <role-name>\n")
		os.Exit(1)
	}

	// Try to find task in global list
	roleName := args[0]
	role, ok := config.Roles[roleName]
	if !ok {
		fmt.Printf("Error: Role not found: %v\n", roleName)
		os.Exit(1)
	}

	logger.Infof("Running role `%v`", roleName)
	for _, host := range config.Hosts.WithRole(roleName) {
		for _, task := range role.Tasks {
			runWithRunner(&task, &host)
		}
	}
}
