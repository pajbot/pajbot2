package modules

import (
	"github.com/pajlada/pajbot2/pkg"
)

type ActionPerformer struct {
	botChannel pkg.BotChannel

	server *server
}

var actionPerformerModuleSpec = &moduleSpec{
	id:   "action_performer",
	name: "Action performer",

	maker: NewActionPerformer,

	enabledByDefault: true,

	priority: 500000,
}

func NewActionPerformer() pkg.Module {
	return &ActionPerformer{
		server: &_server,
	}
}

func (m *ActionPerformer) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel
	return nil
}

func (m *ActionPerformer) Disable() error {
	return nil
}

func (m *ActionPerformer) Spec() pkg.ModuleSpec {
	return actionPerformerModuleSpec
}

func (m *ActionPerformer) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m ActionPerformer) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	return nil
}

func (m ActionPerformer) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return action.Do()
}
