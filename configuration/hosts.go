package configuration

// Host is the definition of a host and its
// specific tasks.
type Host struct {
	Host       string
	User       string
	Password   string
	PrivateKey string `yaml:"private_key"`
	Connection string
	Sudo       bool
	Port       int
	Roles      []string
	Tasks      TaskCollection
}

// HostCollection is a set of hosts.
type HostCollection []Host

// Get finds the host with the given name in the
// collection.
func (h HostCollection) Get(name string) *Host {
	for _, v := range h {
		if v.Host == name {
			return &v
		}
	}

	return nil
}
