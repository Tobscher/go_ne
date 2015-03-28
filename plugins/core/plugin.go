package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RunCommandRegular(name string, arg ...string) error {
	commands := strings.Split(name, " ")
	command := commands[0]
	arguments := append([]string{}, commands[1:]...)
	arguments = append(arguments, arg...)
	cmd := exec.Command(command, arguments...)

	fmt.Printf("Running command %v with args %v \n", command, strings.Join(arguments, " "))

	// only if verbose
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func RunCommandAsSudo(name string, arg ...string) error {
	args := append([]string{name}, arg...)
	return RunCommandRegular("sudo", args...)
}

func RunCommand(sudo bool, name string, arg ...string) error {
	if sudo {
		return RunCommandAsSudo(name, arg...)
	}
	return RunCommandRegular(name, arg...)
}

func LoadConfig(reader io.Reader, v interface{}) {
	bio := bufio.NewReader(reader)
	bytes, _, err := bio.ReadLine()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		log.Fatal(err)
	}
}
