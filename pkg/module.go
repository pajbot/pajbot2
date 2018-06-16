package pkg

type Module interface {
	Name() string
	Register() error
	OnWhisper(source User, message Message) error
	OnMessage(source Channel, user User, message Message, action Action) error
}
