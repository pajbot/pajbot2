package command

// Command xD
type Command struct {
	Trigger  string
	Response string
}

// GetResponse xD
func (command *Command) GetResponse() string {
	// TODO: get $(user) variables and shit
	return command.Response
}
