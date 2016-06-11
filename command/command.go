package command

// Command xD
type Command struct {
	Triggers []string
	Response string
}

/*
IsTriggered returns true if the given string `message` would trigger this command,
otherwise return false
*/
func (command *Command) IsTriggered(message string) bool {
	for _, trigger := range command.Triggers {
		if trigger == message {
			return true
		}
	}
	return false
}

// GetResponse xD
func (command *Command) GetResponse() string {
	// TODO: get $(user) variables and shit
	return command.Response
}
