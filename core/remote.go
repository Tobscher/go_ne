package core

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/tobscher/kiss/configuration"
	"github.com/tobscher/kiss/logging"
)

// Remote describes a runner which runs task
// on a remote system via SSH.
type Remote struct {
	Client *ssh.Client

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

	return &Remote{
		Client: client,
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

func (r *Remote) Prepare() error {
	logger.Debug("Preparing remote machine")

	t := time.Now().Local()
	r.tempDir = fmt.Sprintf("/tmp/kiss/%v", t.Format("20060102150405"))
	cmd := fmt.Sprintf("mkdir -p %v", r.tempDir)
	return r.runCommand(cmd)
}

// Run runs the given task on the remote system
func (r *Remote) Run(t *configuration.Task) error {
	session, err := r.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}
	defer session.Close()

	cmd := "/tmp/kiss/agent"

	if logger.Level > logging.INFO {
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
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

	return session.Wait()
}

// Close closes the SSH connection to the remote system
func (r *Remote) Close() {
	logger.Debug("Tearing down remote machine")

	cmd := fmt.Sprintf("rm -rf %v", r.tempDir)

	r.runCommand(cmd)

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
