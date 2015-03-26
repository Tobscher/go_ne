package core

import "github.com/tobscher/kiss/configuration"

// Args describes the args for the remote call.
type Args struct {
	Environment []string
	Args        []string
	Options     configuration.OptionCollection
}
