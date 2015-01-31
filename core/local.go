package core

import (
	"io"
	"os/exec"

	"log"
	"strings"
)

// ReadBufferSize describes the buffer size to be used
// when reading from stdout/stderr
const ReadBufferSize = 1024

// Local describes a local runner.
type Local struct {
	ChStdOut chan []byte
	ChStdErr chan []byte
}

// NewLocalRunner creates a new runner which runs tasks
// locally.
func NewLocalRunner() (*Local, error) {
	lr := Local{
		ChStdOut: make(chan []byte),
		ChStdErr: make(chan []byte),
	}

	return &lr, nil
}

// Run executes the given task on the local machine.
func (l *Local) Run(task Task) error {
	log.Printf("Running task: %v %v\n", task.Name(), strings.Join(task.Args(), " "))

	cmd := exec.Command(task.Name(), task.Args()...)
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdOut.Close()

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stdErr.Close()

	err = cmd.Start()
	if err != nil {
		return err
	}

	// COULDDO: Might be slicker to spin up two async processes that communicate back
	bufferStdOut := make([]byte, ReadBufferSize)
	bufferStdErr := make([]byte, ReadBufferSize)
	listeningOut := true
	listeningErr := true
	for {
		if listeningOut {
			outBytes, err := stdOut.Read(bufferStdOut)
			if err != nil {
				if err == io.EOF {
					listeningOut = false
				} else {
					return err
				}
			}

			if outBytes > 0 {
				l.ChStdOut <- bufferStdOut[:outBytes]
			}
		}

		if listeningErr {
			errBytes, err := stdOut.Read(bufferStdErr)
			if err != nil {
				if err == io.EOF {
					listeningErr = false
				} else {
					return err
				}
			}

			if errBytes > 0 {
				l.ChStdErr <- bufferStdErr[:errBytes]
			}
		}

		if !listeningOut && !listeningErr { // both complete
			break
		}
	}

	cmd.Wait()
	return nil
}

// Close closes any open connections
func (l *Local) Close() {
}

// LogOutput logs the output received by the channel.
func LogOutput(runner *Local) {
	for {
		select {
		case out := <-runner.ChStdOut:
			log.Printf("%s", out)
		case out := <-runner.ChStdErr:
			log.Printf("%s", out)
		}
	}
}
