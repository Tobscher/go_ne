package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"github.com/tobscher/kiss/configuration"
	"github.com/tobscher/kiss/logging"
)

// Remote describes a runner which runs task
// on a remote system via SSH.
type Remote struct {
	Host       *configuration.Host
	Client     *ssh.Client
	SftpClient *sftp.Client

	tempDir string
}

// NewRemoteRunner creates a new runner which runs
// tasks on a remote system.
func NewRemoteRunner(host *configuration.Host) *Remote {
	return &Remote{
		Host: host,
	}
}

// An SSH connection will be established.
func (r *Remote) Connect() error {
	client, err := createClient(r.Host.User, r.Host.Password, r.Host.Host, strconv.Itoa(r.Host.Port), r.Host.PrivateKey)
	if err != nil {
		return err
	}

	fileClient, err := createFileClient(client)
	if err != nil {
		return err
	}

	r.Client = client
	r.SftpClient = fileClient

	return nil
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

	_, err := r.SftpClient.Lstat(path)
	if err != nil {
		return false
	}

	return true
}

// Run runs the given task on the remote system
func (r *Remote) Run(t *configuration.Task) error {
	if err := r.Prepare(); err != nil {
		return err
	}

	if t.WaitBefore > 0 {
		logger.Infof("Waiting for %v seconds\n", t.WaitBefore)
		time.Sleep(time.Duration(t.WaitBefore) * time.Second)
	}

	if err := r.Execute(t); err != nil {
		return err
	}

	if t.WaitAfter > 0 {
		logger.Infof("Waiting for %v seconds\n", t.WaitAfter)
		time.Sleep(time.Duration(t.WaitAfter) * time.Second)
	}

	if err := r.CleanUp(); err != nil {
		return err
	}

	return nil
}

func (r *Remote) Prepare() error {
	logger.Debug("Preparing remote machine")

	r.Connect()
	defer r.Close()
	r.makeDirectory("/tmp/kiss")

	t := time.Now().Local()
	r.tempDir = fmt.Sprintf("/tmp/kiss/%v", t.Format("20060102150405"))
	if err := r.makeDirectory(r.tempDir); err != nil {
		return err
	}

	return nil
}

func (r *Remote) Execute(t *configuration.Task) error {
	logger.Debug("Executing task...")

	r.Connect()
	defer r.Close()

	session, err := r.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}

	commands := []string{
		"env",
		fmt.Sprintf("KISS_TMP_DIR=%v", r.tempDir),
		"kiss-agent",
	}
	cmd := strings.Join(commands, " ")

	if logger.Level > logging.INFO {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Error while getting stdout pipe: %v", err)
		}

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				logger.Debug(scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				logger.Fatal(fmt.Sprintf("Error while reading from stdout: %v", err))
			}
		}()
	}

	stdErr, err := session.StderrPipe()
	if err != nil {
		return err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err := session.Start(cmd); err != nil {
		return err
	}

	io.WriteString(stdin, t.JSON())
	io.WriteString(stdin, "\n")

	err = session.Wait()
	if err != nil {
		bytes, bufErr := ioutil.ReadAll(stdErr)
		if bufErr != nil {
			return fmt.Errorf("Error reading from stderr: %v", bufErr)
		}

		return fmt.Errorf("%v: %v", err, string(bytes))
	}

	return nil
}

// Close closes the SSH connection to the remote system
func (r *Remote) Close() {
	logger.Debug("Closing connection to remote machine")
	r.SftpClient.Close()
	r.Client.Close()
}

func (r *Remote) CleanUp() error {
	logger.Debug("Cleaning up remote machine")
	r.Connect()
	defer r.Close()

	return r.remove(r.tempDir, false)
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

	logger.Debugf("Connecting to %v@%v", username, remoteServer)
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
