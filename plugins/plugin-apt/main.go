package main

import (
	"os"

	"github.com/tobscher/kiss/plugins/core"
)

func main() {
	var options Options
	plugin.LoadConfig(os.Stdin, &options)

	apt := NewApt(options)
	result := apt.Run()

	os.Exit(result)
}
