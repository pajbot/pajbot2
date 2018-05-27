package pkg

type Sender interface {
	Say(Channel, string)
	Timeout(Channel, User, int, string)
}
