package core

import (
	"fmt"
	"os/exec"
	"path"
)

var (
	compiler = "goxc"
)

func compileDirectory(directory string, targetOs string, targetArch string) (*string, error) {
	osArgument := fmt.Sprintf("-os=%v", targetOs)
	archArgument := fmt.Sprintf("-arch=%v", targetArch)

	cmd := exec.Command(compiler, osArgument, archArgument, "-d=./tmp", "-tasks=xc")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	cmd.Dir = directory

	go DebugLines(stdout)
	go DebugLines(stderr)

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	directoryName := path.Base(directory)
	file := fmt.Sprintf("%v/tmp/snapshot/%v_%v/%v", directory, targetOs, targetArch, directoryName)
	return &file, nil
}
