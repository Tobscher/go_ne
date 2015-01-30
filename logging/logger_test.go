package logging

import "testing"

func TestLogger(t *testing.T) {
	logger := GetLogger("test")
	logger.SetLevel(TRACE)

	logger.Trace("Trace")
	logger.Tracef("Trace - %v", "foo")

	logger.Debug("Debug")
	logger.Debugf("Debug - %v", "foo")

	logger.Info("Info")
	logger.Infof("Info - %v", "foo")

	logger.Warn("Warn")
	logger.Warnf("Warn - %v", "foo")

	logger.Error("Error")
	logger.Errorf("Error - %v", "foo")

	logger.Fatal("Fatal")
	logger.Fatalf("Fatal - %v", "foo")
}
