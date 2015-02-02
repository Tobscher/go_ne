package configuration

type RoleCollection map[string]Role

type Role struct {
	With  []string
	Tasks TaskCollection
}
