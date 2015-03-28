package commands

import "github.com/tobscher/kiss/logging"

var (
	logger     = logging.GetLogger("kiss")
	configFile string
	verbose    bool
	trace      bool
)
