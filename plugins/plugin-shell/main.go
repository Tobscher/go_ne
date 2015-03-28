package main

import (
	"os"

	"github.com/tobscher/kiss/plugins/core"
)

func main() {
	var options Options
	plugin.LoadConfig(os.Stdin, &options)

	shell := NewShell(options)
	result := shell.Run()

	os.Exit(result)
}
