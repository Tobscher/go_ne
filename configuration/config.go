package configuration

import (
	"io/ioutil"

	"github.com/tobscher/kiss/logging"
	"gopkg.in/yaml.v2"
)

var logger = logging.GetLogger("kiss")

// Configuration holds information about
// variables, hosts and tasks.
type Configuration struct {
	Vars  map[string]string
	Hosts HostCollection
	Tasks TaskCollection
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
