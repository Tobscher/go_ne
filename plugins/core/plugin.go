package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func RunCommand(name string, sudo bool) error {
	bytes := []byte(name)
	tempFile, err := ioutil.TempFile(os.Getenv("KISS_TMP_DIR"), "shell-")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	fmt.Printf("Created temporary file `%v`\n", tempFile.Name())

	tempFile.Write(bytes)

	var cmd *exec.Cmd
	if sudo {
		cmd = exec.Command(sudoCommand, shellCommand, tempFile.Name())
	} else {
		cmd = exec.Command(shellCommand, tempFile.Name())
	}

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
