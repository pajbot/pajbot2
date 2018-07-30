package modules

import "github.com/pajlada/pajbot2/pkg"

var _ pkg.Module = &Test{}

type Test struct {
	server *server
}

func NewTest() *Test {
	return &Test{
		server: &_server,
	}
}

func (m *Test) Register() error {
	return nil
}

func (m *Test) Name() string {
	return "Test"
}

func (m Test) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m Test) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return nil
}
