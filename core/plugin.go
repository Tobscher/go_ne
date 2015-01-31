package core

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/tobscher/go_ne/plugins/shared"
)

const (
	pluginPrefix = "plugin"
	maxAttempts  = 5 // to connect to plugin
)

// PluginCache stores loaded plugins
type PluginCache struct {
	sync.Mutex
	cache map[string]*Plugin
}

var (
	startPort     = 8000
	loadedPlugins = PluginCache{
		cache: make(map[string]*Plugin),
	}
)

// Plugin stores information about the plugin and how to
// connect to it.
type Plugin struct {
	information *PluginInformation
	client      *rpc.Client
}

// PluginInformation stores the details about the plugin
// such as the host, port and the underlying command
type PluginInformation struct {
	Host string
	Port string
	Cmd  *exec.Cmd
}

// Address returns the full address and port, e.g. localhost:8001
func (p *PluginInformation) Address() string {
	return fmt.Sprintf("%v:%v", p.Host, p.Port)
}

// StartPlugin starts the plugin with the given name.
// This will try to boot an application called `plugin-<plugin-name>`
//
// This method will return an error when the plugin can not be found
// or the plugin exits with an exit code other than 0.
func StartPlugin(name string) (*Plugin, error) {
	command := fmt.Sprintf("%v-%v", pluginPrefix, name)
	host := "localhost"
	port := nextAvailblePort()

	logger.Debugf("Starting plugin `%v` on port %v", name, port)

	// Pass host and port to plugin
	cmd := exec.Command(command,
		fmt.Sprintf("-host=%v", host),
		fmt.Sprintf("-port=%v", port),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the plugin
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	info := &PluginInformation{
		Host: host,
		Port: port,
		Cmd:  cmd,
	}

	// Try to establish connection to plugin
	var conn net.Conn
	for i := 1; i <= maxAttempts; i++ {
		logger.Tracef("Attempt %v to connect to plugin...", i)

		conn, err = net.Dial("tcp", info.Address())
		if err != nil {
			time.Sleep(100 * time.Millisecond)

			if i == maxAttempts {
				cmd.Process.Kill()
				return nil, errors.New("Could not connect to plugin.")
			}

			continue
		}

		break
	}

	client := jsonrpc.NewClient(conn)

	plugin := &Plugin{
		information: info,
		client:      client,
	}

	loadedPlugins.Lock()
	loadedPlugins.cache[name] = plugin
	loadedPlugins.Unlock()

	return plugin, nil
}

// GetPlugin returns a loaded plugin. If the plugin has
// not been loaded yet, it will load it.
func GetPlugin(name string) (*Plugin, error) {
	var val *Plugin
	var ok bool
	var err error

	val, ok = loadedPlugins.cache[name]
	if !ok {
		val, err = StartPlugin(name)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}

// GetCommands asks the plugin (via RPC) which commands should
// be executed on the remote system.
func (p *Plugin) GetCommands(args shared.Args) ([]*Command, error) {
	var reply shared.Response
	var commands []*Command

	err := p.client.Call("Command.Execute", args, &reply)
	if err != nil {
		return nil, err
	}

	for _, value := range reply.Commands {
		command := &Command{
			name: value.Name,
			args: value.Args,
		}

		commands = append(commands, command)
	}

	return commands, nil
}

func nextAvailblePort() string {
	startPort++
	return strconv.Itoa(startPort)
}

// StopAllPlugins will stop all plugins which are currently running.
//
// BUG(Tobscher) Send signal to gracefully shutdown the plugin
func StopAllPlugins() {
	loadedPlugins.Lock()
	defer loadedPlugins.Unlock()

	for k, v := range loadedPlugins.cache {
		logger.Debugf("Stopping plugin: %v", k)
		if err := v.information.Cmd.Process.Kill(); err != nil {
			logger.Warn(err.Error())
		}
	}

	loadedPlugins.cache = make(map[string]*Plugin)
	startPort = 8000
}
