package main

import (
	"fmt"
	"log"
	"strings"

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
	return plugin.RunCommand("apt-get update -qq", sudo)
}

func installPackages(sudo bool, packages []string) error {
	packageList := strings.Join(packages, " ")
	return plugin.RunCommand(fmt.Sprintf("apt-get install -y %v", packageList), sudo)
}
