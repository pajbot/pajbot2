package modules

import "github.com/pajlada/pajbot2/pkg"

var _ pkg.Module = &MessageLengthLimit{}

type MessageLengthLimit struct {
	botChannel pkg.BotChannel

	server *server
}

func newMessageLengthLimit() pkg.Module {
	return &MessageLengthLimit{
		server: &_server,
	}
}

var messageLengthLimitSpec = moduleSpec{
	id:    "message_length_limit",
	name:  "Message length limit",
	maker: newMessageLengthLimit,
}

func (m *MessageLengthLimit) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	return nil
}

func (m *MessageLengthLimit) Disable() error {
	return nil
}

func (m *MessageLengthLimit) Spec() pkg.ModuleSpec {
	return &messageLengthLimitSpec
}

func (m *MessageLengthLimit) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m MessageLengthLimit) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	return nil
}

func (m MessageLengthLimit) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	return nil

	if user.GetName() == "gazatu2" {
		return nil
	}

	if user.GetName() == "supibot" {
		return nil
	}

	messageLength := len(message.GetText())
	if messageLength > 140 {
		if messageLength > 420 {
			action.Set(pkg.Timeout{600, "Your message is way too long"})
			return nil
		}

		action.Set(pkg.Timeout{300, "Your message is too long, shorten it"})
		return nil
	}

	return nil
}
