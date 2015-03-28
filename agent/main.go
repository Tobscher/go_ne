package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/tobscher/kiss/configuration"
)

func main() {
	bio := bufio.NewReader(os.Stdin)
	bytes, _, err := bio.ReadLine()
	if err != nil {
		exitWithError(err)
	}

	var task configuration.Task
	err = json.Unmarshal(bytes, &task)
	if err != nil {
		exitWithError(err)
	}

	err = RunTask(&task)
	if err != nil {
		exitWithError(err)
	}

	os.Exit(0)
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
