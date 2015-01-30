package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new command to output the current version
func NewVersionCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print the version number of %v", name),
		Long:  fmt.Sprintf(`Print the version number of %v`, name),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v %v\n", name, version)
		},
	}

	return command
}
