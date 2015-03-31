package core

import "fmt"

var (
	installDirectory = ".kiss"
	agentName        = "agent"
	agent            = fmt.Sprintf("%v/bin/%v", installDirectory, agentName)
	pluginDirectory  = fmt.Sprintf("%v/bin", installDirectory)
	pluginPrefix     = "plugin"
)
