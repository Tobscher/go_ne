package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	Config     *configuration.Configuration
	Facts      Facts
	Client     *ssh.Client
	SftpClient *sftp.Client

	tempDir string
}

// NewRemoteRunner creates a new runner which runs
// tasks on a remote system.
func NewRemoteRunner(host *configuration.Host, config *configuration.Configuration) *Remote {
	return &Remote{
		Host:   host,
		Config: config,
	}
}

// RunCommand runs an abritrary command on the remote system.
func (r *Remote) GatherFacts() (Facts, error) {
	if err := r.Connect(); err != nil {
		return nil, err
	}

	defer r.Close()

	session, err := r.Client.NewSession()
	if err != nil {
		return nil, errors.New("Failed to create session: " + err.Error())
	}
	defer session.Close()

	logger.Debugf("Gathering facts")

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := session.Start(gather); err != nil {
		return nil, err
	}

	if err = session.Wait(); err != nil {
		return nil, err
	}

	facts := make(Facts)

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " = ")
		facts[parts[0]] = parts[1]
	}

	return facts, nil
}

// Run runs the given task on the remote system.
func (r *Remote) Run(t *configuration.Task) error {
	if err := r.Prepare(t); err != nil {
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

func (r *Remote) BeforeAll(tasks configuration.TaskCollection) error {
	// Gather facts
	facts, err := r.GatherFacts()
	if err != nil {
		return err
	}

	r.Facts = facts

	// Check for agent
	if r.Config.Agent.Force {
		r.remove(agent, false)
	}

	if !r.fileExists(agent) {
		logger.Warnf("Agent is not installed on the remote system: %v", agent)

		// When configured to use local file
		if r.Config.Agent.Path != "" {
			logger.Infof("Compiling agent for %v/%v", facts.OS(), facts.Arch())

			file, err := compileDirectory(r.Config.Agent.Path, facts.OS(), facts.Arch())
			if err != nil {
				return err
			}

			if err = r.uploadFile(*file, agent); err != nil {
				return err
			}

			if err = r.chmod(agent, os.FileMode(0755)); err != nil {
				return err
			}
		} else {
			downloadUrl := bintrayDownloadUrl(agentName, facts.OS(), facts.Arch())
			r.downloadFile(downloadUrl)
		}
	}

	// Check for plugins
	for _, pluginName := range tasks.UniquePluginNames() {
		plugin := fmt.Sprintf("%v/%v-%v", pluginDirectory, pluginPrefix, pluginName)
		pluginConfig := r.Config.Plugins[pluginName]

		if pluginConfig.Force {
			r.remove(plugin, false)
		}

		if !r.fileExists(plugin) {
			logger.Warnf("Plugin is not installed on the remote system: %v", plugin)

			if pluginConfig.Path != "" {
				logger.Infof("Compiling %v-%v for %v/%v", pluginPrefix, pluginName, r.Facts.OS(), r.Facts.Arch())
				file, err := compileDirectory(r.Config.Plugins[pluginName].Path, r.Facts.OS(), r.Facts.Arch())
				if err != nil {
					return err
				}

				if err = r.uploadFile(*file, plugin); err != nil {
					return err
				}

				if err = r.chmod(plugin, os.FileMode(0755)); err != nil {
					return err
				}
			} else {
				downloadUrl := bintrayDownloadUrl(fmt.Sprintf("%v-%v", pluginPrefix, pluginName), facts.OS(), facts.Arch())
				r.downloadFile(downloadUrl)
			}
		}
	}

	return nil
}

// Prepare prepares the remote system so it can run plugins. This will do the following:
// * Create temp directory /tmp/kiss/<current_datetime>
func (r *Remote) Prepare(task *configuration.Task) error {
	logger.Debug("Preparing remote machine")

	if err := r.Connect(); err != nil {
		return err
	}
	defer r.Close()
	r.makeDirectory("/tmp/kiss")

	t := time.Now().Local()
	r.tempDir = fmt.Sprintf("/tmp/kiss/%v", t.Format("20060102150405"))
	if err := r.makeDirectory(r.tempDir); err != nil {
		return err
	}

	return nil
}

// Execute executes the task on the remote system via the runner.
func (r *Remote) Execute(t *configuration.Task) error {
	logger.Debug("Executing task...")

	if err := r.Connect(); err != nil {
		return err
	}
	defer r.Close()

	session, err := r.Client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}

	commands := []string{
		"env",
		fmt.Sprintf("KISS_TMP_DIR=%v", r.tempDir),
		agent,
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
		return fmt.Errorf("Command error: %v", err)
	}

	io.WriteString(stdin, t.JSON())
	io.WriteString(stdin, "\n")

	err = session.Wait()
	if err != nil {
		bytes, bufErr := ioutil.ReadAll(stdErr)
		if bufErr != nil {
			return fmt.Errorf("Error reading from stderr: %v", bufErr)
		}

		return fmt.Errorf("%v: %v", err.Error(), string(bytes))
	}

	return nil
}

// CleanUp removes the temporary directory on the remote system.
func (r *Remote) CleanUp() error {
	logger.Debug("Cleaning up remote machine")

	if err := r.Connect(); err != nil {
		return err
	}
	defer r.Close()

	return r.remove(r.tempDir, false)
}

// AfterAll runs post task scripts.
func (r *Remote) AfterAll() error {
	return nil
}
