package main

import (
	"log"
	"os"

	"github.com/tobscher/kiss/plugins/core"
)

func main() {
	f, err := os.OpenFile("/tmp/shell.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	log.SetOutput(f)

	var options Options
	plugin.LoadConfig(os.Stdin, &options)

	shell := NewShell(options)
	result := shell.Run()

	os.Exit(result)
}
