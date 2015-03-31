package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Connct connects to the remote system via SSH.
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

// Close closes the SSH connection to the remote system
func (r *Remote) Close() {
	logger.Debug("Closing connection to remote machine")
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
