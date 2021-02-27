package data

//Command contains actions and data relating to those actions like, help, roles, etc.
type Command struct {
	Action func()
	Help   string
	Roles  []string
}
