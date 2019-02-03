package pkg

type CustomCommand interface {
	Trigger(BotChannel, []string, Channel, User, Message, Action)
}
