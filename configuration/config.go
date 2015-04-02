package configuration

import (
	"io/ioutil"

	"github.com/tobscher/kiss/logging"
	"gopkg.in/yaml.v2"
)

var logger = logging.GetLogger("kiss")

type Agent struct {
	Path  string
	Force bool
}

type PluginConfig struct {
	Path  string
	Force bool
}

// Configuration holds information about
// variables, hosts and tasks.
type Configuration struct {
	Vars    map[string]string
	Hosts   HostCollection
	Tasks   TaskCollection
	Roles   RoleCollection
	Agent   Agent
	Plugins map[string]PluginConfig
}

// Load loads configuration from the given path.
func Load(file string) *Configuration {
	rawYaml, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Warnf("Could not load configuration from `%v`.", file)
		return nil
	}

	var c Configuration
	yaml.Unmarshal(rawYaml, &c)

	return &c
}
