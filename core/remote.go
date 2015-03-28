package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/tobscher/kiss/configuration"
	"github.com/tobscher/kiss/logging"
)

// Remote describes a runner which runs task
// on a remote system via SSH.
type Remote struct {
	Client     *ssh.Client
	SftpClient *sftp.Client

	tempDir string
}

// NewRemoteRunner creates a new runner which runs
// tasks on a remote system.
//
// An SSH connection will be establishe.
func NewRemoteRunner(host *configuration.Host) (*Remote, error) {
	client, err := createClient(host.User, host.Password, host.Host, strconv.Itoa(host.Port), host.PrivateKey)
	if err != nil {
		return nil, err
	}

	fileClient, err := createFileClient(client)
	if err != nil {
		return nil, err
	}

	return &Remote{
		Client:     client,
		SftpClient: fileClient,
	}, nil
}

func (r *Remote) runCommand(cmd string) error {
	session, err := r.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}
	defer session.Close()

	logger.Debugf("Executing `%v`", cmd)
	if err := session.Start(cmd); err != nil {
		return err
	}

	if err = session.Wait(); err != nil {
		return err
	}

	return nil
}

func (r *Remote) uploadFile(from string, to string) error {
	logger.Tracef("Uploading file from `%v` to `%v`\n", from, to)

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

func (r *Remote) remove(path string) error {
	logger.Tracef("Removing file/directory `%v`\n", path)
	err := r.SftpClient.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (r *Remote) fileExists(path string) bool {
	logger.Tracef("Checking if file exists `%v`\n", path)

	_, err := r.SftpClient.Lstat(path)
	if err != nil {
		return false
	}

	return true
}

func (r *Remote) Prepare() error {
	logger.Debug("Preparing remote machine")

	t := time.Now().Local()
	r.tempDir = fmt.Sprintf("/tmp/kiss/%v", t.Format("20060102150405"))
	return r.makeDirectory(r.tempDir)
}

// Run runs the given task on the remote system
func (r *Remote) Run(t *configuration.Task) error {
	session, err := r.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}
	defer session.Close()

	cmd := "kiss-agent"

	if logger.Level > logging.INFO {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return err
		}

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				logger.Debug(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				logger.Fatal(err.Error())
			}
		}()
	}

	stdErr, err := session.StderrPipe()
	if err != nil {
		return err
	}

	stdin, err := session.StdinPipe()
	if err := session.Start(cmd); err != nil {
		return err
	}

	if err != nil {
		logger.Fatal(err.Error())
	}

	io.WriteString(stdin, t.JSON())
	io.WriteString(stdin, "\n")

	err = session.Wait()
	if err != nil {
		bytes, bufErr := ioutil.ReadAll(stdErr)
		if bufErr != nil {
			return bufErr
		}

		return fmt.Errorf("%v: %v", err, string(bytes))
	}

	return nil
}

// Close closes the SSH connection to the remote system
func (r *Remote) Close() {
	logger.Debug("Tearing down remote machine")

	r.remove(r.tempDir)
	r.SftpClient.Close()
	r.Client.Close()
}

func createClient(username, password, host, port, key string) (*ssh.Client, error) {
	authMethods := []ssh.AuthMethod{}

	if len(password) > 0 {
		authMethods = append(authMethods, ssh.Password(password))
	}

	if len(key) > 0 {
		priv, err := loadKey(key)
		if err != nil {
			log.Println(err)
		} else {
			signers, err := ssh.NewSignerFromKey(priv)
			if err != nil {
				log.Println(err)
			} else {
				authMethods = append(authMethods, ssh.PublicKeys(signers))
			}
		}
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: authMethods,
	}

	if len(port) == 0 {
		port = "22"
	}

	remoteServer := fmt.Sprintf("%v:%v", host, port)

	logger.Infof("Connecting to %v@%v", username, remoteServer)
	client, err := ssh.Dial("tcp", remoteServer, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createFileClient(client *ssh.Client) (*sftp.Client, error) {
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}

	return sftp, nil
}

func loadKey(file string) (interface{}, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParseRawPrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return key, nil
}
