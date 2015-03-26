package core

import "github.com/tobscher/kiss/configuration"

// Runner is the interface which describes objects
// that can execute tasks on a system.
type Runner interface {
	Prepare() error
	Run(*configuration.Task) error
	Close()
}
