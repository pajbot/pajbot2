package pkg

type CustomCommand interface {
	Trigger(BotChannel, []string, Channel, User, Message, Action)
}

type CommandInfo struct {
	Name        string
	Description string

	Maker func() CustomCommand2 `json:"-"`
}

type CustomCommand2 interface {
	Trigger(BotChannel, []string, User, Message, Action)
	HasCooldown(User) bool
	AddCooldown(User)
}

type CommandMatcher interface {
	Register(aliases []string, command interface{}) interface{}
	Deregister(command interface{})
	Match(text string) (interface{}, []string)
}

type CommandsManager interface {
	CommandMatcher
	OnMessage(bot BotChannel, user User, message Message, action Action) error
}
