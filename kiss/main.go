package main

func main() {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewRunCommand())

	rootCmd.Execute()
}
