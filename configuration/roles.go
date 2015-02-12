package configuration

type RoleCollection map[string]Role

type Role struct {
	With        []string
	Description string
	Tasks       TaskCollection
}
