package configuration

import (
	"encoding/json"
)

// OptionCollection is a map from option name
// to option value
type OptionCollection map[string]interface{}

// Plugin has information about the plugin and
// its options.
type Plugin struct {
	Options OptionCollection `yaml:",inline"`
}

// Task is a definition of a command to be run on the
// remote server.
type Task struct {
	Task        string
	Description string
	Plugin      map[string]Plugin `yaml:",inline"`
}

// TaskCollection defines a set of tasks.
type TaskCollection []Task

// Get returns the task with the given name in
// the collection
func (t TaskCollection) Get(name string) *Task {
	for _, v := range t {
		if v.Task == name {
			return &v
		}
	}

	return nil
}

func (t *Task) JSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return string(bytes)
}
