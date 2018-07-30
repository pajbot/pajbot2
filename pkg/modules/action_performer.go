package modules

import "github.com/pajlada/pajbot2/pkg"

var _ pkg.Module = &ActionPerformer{}

type ActionPerformer struct {
	server *server
}

func NewActionPerformer() *ActionPerformer {
	return &ActionPerformer{
		server: &_server,
	}
}

func (m *ActionPerformer) Register() error {
	return nil
}

func (m *ActionPerformer) Name() string {
	return "ActionPerformer"
}

func (m ActionPerformer) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m ActionPerformer) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return action.Do()
}
