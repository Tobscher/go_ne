package main

import "github.com/tobscher/go_ne/commands"

func main() {
	rootCmd := commands.NewRootCommand()
	rootCmd.AddCommand(commands.NewVersionCommand(name, version))
	rootCmd.AddCommand(commands.NewRunCommand())
	rootCmd.AddCommand(commands.NewTasksCommand())
	rootCmd.AddCommand(commands.NewHostsCommand())

	rootCmd.Execute()
}
