package pkg

type Module interface {
	Name() string
	Register() error
	OnWhisper(bot Sender, source User, message Message) error
	OnMessage(bot Sender, source Channel, user User, message Message, action Action) error
}
