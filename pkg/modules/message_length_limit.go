package modules

import "github.com/pajlada/pajbot2/pkg"

var _ pkg.Module = &MessageLengthLimit{}

type MessageLengthLimit struct {
	server *server
}

func NewMessageLengthLimit() *MessageLengthLimit {
	return &MessageLengthLimit{
		server: &_server,
	}
}

func (m *MessageLengthLimit) Register() error {
	return nil
}

func (m *MessageLengthLimit) Name() string {
	return "MessageLengthLimit"
}

func (m MessageLengthLimit) OnWhisper(user pkg.User, message pkg.Message) error {
	return nil
}

func (m MessageLengthLimit) OnMessage(channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if channel.GetChannel() != "forsen" {
		return nil
	}

	if user.GetName() == "gazatu2" {
		return nil
	}

	if user.GetName() == "supibot" {
		return nil
	}

	messageLength := len(message.GetText())
	if messageLength > 140 {
		if messageLength > 420 {
			action.SetTimeout(600, "Your message is way too long")
			return nil
		}

		action.SetTimeout(300, "Your message is too long, shorten it")
		return nil
	}

	return nil
}
