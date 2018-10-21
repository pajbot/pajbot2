package pkg

type CustomCommand interface {
	Trigger(Sender, BotChannel, []string, Channel, User, Message, Action)
}
