package core

import "github.com/tobscher/go_ne/logging"

var logger = logging.GetLogger("core")

// SetLogLevel sets the log level for the core logger
func SetLogLevel(level logging.Level) {
	logger.SetLevel(level)
}
