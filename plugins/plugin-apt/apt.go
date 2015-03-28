package main

import (
	"fmt"
	"log"

	"github.com/tobscher/kiss/plugins/core"
)

type Apt struct {
	options Options
}

func NewApt(options Options) *Apt {
	apt := Apt{
		options: options,
	}

	return &apt
}

func (a *Apt) Run() int {
	if a.options.Update {
		fmt.Println("-- Updating cache...")
		if err := updateCache(a.options.Sudo); err != nil {
			log.Println(err)
			return 1
		}
	}

	if len(a.options.Packages) > 0 {
		fmt.Println("-- Installing packages...")
		if err := installPackages(a.options.Sudo, a.options.Packages); err != nil {
			log.Println(err)
			return 1
		}
	}

	return 0
}

func updateCache(sudo bool) error {
	return plugin.RunCommand(sudo, "apt-get", "update", "-qq")
}

func installPackages(sudo bool, packages []string) error {
	args := append([]string{"install", "-y"}, packages...)
	return plugin.RunCommand(sudo, "apt-get", args...)
}
