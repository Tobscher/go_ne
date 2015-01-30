package main

import "github.com/spf13/cobra"

// NewRootCommand creates the root command for the CLI.
func NewRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:  "kiss",
		Long: "kiss is a plugin-based server automation and deployment tool.",
	}

	return command
}
