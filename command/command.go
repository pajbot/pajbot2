package command

// Command is the shared interface for all commands
type Command interface {
	IsTriggered(t string, fullMessage []string, index int) (bool, Command)
	Run() string
}
