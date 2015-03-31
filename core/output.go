package core

import (
	"bufio"
	"io"
)

func DebugLines(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logger.Debug(scanner.Text())
	}
}
