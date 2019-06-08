package modules

import "github.com/pajbot/pajbot2/pkg"

func init() {
	Register("message_length_limit", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "message_length_limit",
			name:  "Message length limit",
			maker: newMessageLengthLimit,
		}
	})
}

var _ pkg.Module = &MessageLengthLimit{}

type MessageLengthLimit struct {
	base
}

func newMessageLengthLimit(b base) pkg.Module {
	return &MessageLengthLimit{
		base: b,
	}
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
