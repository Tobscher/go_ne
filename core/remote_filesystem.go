package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (r *Remote) uploadFile(from string, to string) error {
	logger.Debugf("Uploading file from `%v` to `%v`\n", from, to)

	if err := r.Connect(); err != nil {
		return err
	}

	defer r.Close()

	directory := filepath.Dir(to)
	directories := strings.Split(directory, "/")

	var currentDir []string
	for _, value := range directories {
		if value == "" {
			continue
		}

		currentDir = append(currentDir, value)
		err := r.makeDirectory(strings.Join(currentDir, "/"))
		if err != nil {
			logger.Warn(err.Error())
		}
	}

	f, err := r.SftpClient.Create(to)
	if err != nil {
		return err
	}
	logger.Trace("File created on remote system.")

	file, err := os.Open(from)
	if err != nil {
		return err
	}
	logger.Trace("File opened for read on host system.")

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	logger.Trace("Contents of file read.")

	if _, err := f.Write(bytes); err != nil {
		return err
	}
	logger.Trace("Bytes written to remote file.")

	_, err = r.SftpClient.Lstat(to)
	if err != nil {
		return err
	}
	logger.Trace("File verified.")

	return nil
}

func (r *Remote) downloadFile(file string) error {
	logger.Debugf("Downloading file: %v", file)

	if err := r.Connect(); err != nil {
		return err
	}
	defer r.Close()

	session, err := r.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}

	if err := session.Start(fmt.Sprintf("curl -L %v | tar -xz --strip-components=1 -C .kiss/bin", file)); err != nil {
		return fmt.Errorf("Command error: %v", err)
	}

	err = session.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (r *Remote) makeDirectory(path string) error {
	logger.Tracef("Creating directory `%v`\n", path)

	err := r.SftpClient.Mkdir(path)
	if err != nil {
		return err
	}

	return nil
}

func (r *Remote) remove(path string, checkError bool) error {
	logger.Tracef("Removing file/directory `%v`\n", path)
	err := r.SftpClient.Remove(path)
	if err != nil {
		if checkError {
			return err
		}
	}

	return nil
}

func (r *Remote) fileExists(path string) bool {
	logger.Tracef("Checking if file exists `%v`\n", path)

	if err := r.Connect(); err != nil {
		return false
	}
	defer r.Close()

	_, err := r.SftpClient.Lstat(path)
	if err != nil {
		return false
	}

	return true
}

func (r *Remote) chmod(path string, mode os.FileMode) error {
	logger.Tracef("Changing file mode of `%v` to `%v`\n", path, mode)

	if err := r.Connect(); err != nil {
		return err
	}
	defer r.Close()

	return r.SftpClient.Chmod(path, mode)
}
