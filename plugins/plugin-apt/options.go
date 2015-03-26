package main

type Options struct {
	Update   bool     `yaml:"update"`
	Packages []string `yaml:"packages"`
	Sudo     bool     `yaml:"sudo"`
}
