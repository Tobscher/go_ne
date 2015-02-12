package main

import "github.com/tobscher/kiss/plugins/core"

func main() {
	plugin.Register(NewCommand())
	plugin.Serve()
}
