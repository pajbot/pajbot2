package pkg

type CustomCommand interface {
	Trigger(Sender, []string, Channel, User, Message, Action)
}
