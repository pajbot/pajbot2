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

type CommandsManager interface {
	Register(aliases []string, command CustomCommand2) CustomCommand2
	Deregister(command CustomCommand2)
	OnMessage(bot BotChannel, user User, message Message, action Action) error
}
