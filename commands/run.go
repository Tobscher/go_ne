package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tobscher/go_ne/configuration"
	"github.com/tobscher/go_ne/core"
	"github.com/tobscher/go_ne/logging"
)

var (
	configFile string
	verbose    bool
)

// NewRunCommand creates a new command to run a task.
func NewRunCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "run",
		Short: "Run a task",
		Long:  "Run a task on the defined remote system",
		Run:   runRun,
	}
	command.Flags().StringVar(&configFile, "config", ".kiss.yml", "path to config file")
	command.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	return command
}

func runRun(cmd *cobra.Command, args []string) {
	if verbose {
		logger.SetLevel(logging.DEBUG)
		core.SetLogLevel(logging.DEBUG)
	}

	config := configuration.Load(configFile)

	if len(args) < 1 {
		fmt.Printf("Error: Expected task-name: kiss run <task-name>\n")
		os.Exit(1)
	}

	// Try to find task in global list
	taskName := args[0]
	task := config.Tasks.Get(taskName)
	if task != nil {
		// Can't omit host name for global tasks
		if len(args) < 2 {
			fmt.Printf("Error: Global tasks require host-name: kiss run %v <host-name>\n", taskName)
			os.Exit(1)
		}

		hostName := args[1]
		host := config.Hosts.Get(hostName)
		if host == nil {
			log.Fatalf("Host not found: %v\n", hostName)
		}

		runWithRunner(task, host)
	} else {
		for _, host := range config.Hosts {
			task := host.Tasks.Get(taskName)
			if task != nil {
				runWithRunner(task, &host)
			}
		}
	}

	fmt.Printf("Error: Task not found: %v\n", taskName)
	os.Exit(1)
}

func runWithRunner(task *configuration.Task, host *configuration.Host) {
	logger.Infof("Selecting host `%v`", host.Host)

	runner, err := core.NewRemoteRunner(host)
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	logger.Infof("Executing task `%v`", task.Task)
	err = core.RunTask(runner, task)
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	} else {
		logger.Info("Tasks completed successfully")
	}
	os.Exit(0)
}
