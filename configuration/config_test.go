package configuration

import "testing"

func TestConfigVars(t *testing.T) {
	config := Load("./fixtures/example.yml")

	expected := 2
	vars := len(config.Vars)
	if vars != expected {
		t.Errorf("Expected config.Vars to be %v got %v", expected, vars)
	}

	if config.Vars["application"] != "go-example-app" {
		t.Error("Vars not loaded properly.")
	}

	if config.Vars["version"] != "0.1" {
		t.Error("Vars not loaded properly.")
	}
}

func TestHosts(t *testing.T) {
	config := Load("./fixtures/example.yml")

	expected := 2
	servers := len(config.Hosts)
	if servers != expected {
		t.Errorf("Expected config.Hosts to be %v got %v", expected, servers)
	}

	localhost := config.Hosts.Get("localhost")
	if localhost.Connection != "local" {
		t.Error("Expected connection to be local.")
	}

	if localhost.User != "foo" {
		t.Error("Expected user to be foo.")
	}

	if localhost.Password != "secret" {
		t.Error("Expected password to be secret.")
	}

	if len(localhost.Roles) != 2 {
		t.Error("Expected server localhost to have 2 roles.")
	}

	if localhost.Roles[0] != "docker" {
		t.Error("Expected server localhost to have role docker")
	}

	if localhost.Roles[1] != "web" {
		t.Error("Expected server localhost to have role web")
	}

	if len(localhost.Tasks) != 2 {
		t.Errorf("Expected localhost to have 2 tasks got %v.", len(localhost.Tasks))
	}

	vagrant := config.Hosts.Get("vagrant")
	if vagrant.Connection != "" {
		t.Error("Expected connection to be empty.")
	}

	if !vagrant.Sudo {
		t.Error("Expected vagrant to require sudo.")
	}

	if vagrant.User != "vagrant" {
		t.Error("Expected user to be vagrant.")
	}

	if vagrant.PrivateKey != "$HOME/.ssh/vagrant" {
		t.Errorf("Expected user to be vagrant got %v.", vagrant.PrivateKey)
	}

	if vagrant.Port != 2222 {
		t.Error("Expected port to be 2222.")
	}

	if len(vagrant.Roles) != 2 {
		t.Error("Expected server vagrant to have 2 roles.")
	}

	if vagrant.Roles[0] != "docker" {
		t.Error("Expected server localhost to have role docker")
	}

	if vagrant.Roles[1] != "db" {
		t.Error("Expected server localhost to have role db")
	}
}

func TestTasks(t *testing.T) {
	config := Load("./fixtures/example.yml")

	if len(config.Tasks) != 2 {
		t.Errorf("Expected 2 tasks got %v.", len(config.Tasks))
	}

	env := config.Tasks.Get("Environment")
	shell, _ := env.Plugin["shell"]
	if len(shell.Options) != 1 {
		t.Error("Expected plugin shell to have a command option.")
	}

	if shell.Options["command"] != "env" {
		t.Error("Expected command to be env.")
	}

	gov := config.Tasks.Get("Go version")
	shell, _ = gov.Plugin["shell"]
	if len(shell.Options) != 1 {
		t.Error("Expected plugin shell to have a command option.")
	}

	if shell.Options["command"] != "go version" {
		t.Error("Expected command to be go version.")
	}
}
