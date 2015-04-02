package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	compiler       = "goxc"
	defaultVersion = "snapshot"
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

	version := goxcVersion(directory)

	directoryName := path.Base(directory)
	file := fmt.Sprintf("%v/tmp/%v/%v_%v/%v", directory, version, targetOs, targetArch, directoryName)
	return &file, nil
}

func goxcVersion(directory string) string {
	filePath := strings.Join([]string{directory, ".goxc.json"}, "/")
	file, err := os.Open(filePath)
	if err != nil {
		return defaultVersion
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return defaultVersion
	}

	var config struct {
		PackageVersion string
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return defaultVersion
	}

	return config.PackageVersion
}
