package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	Register("message_length_limit", func() pkg.ModuleSpec {
		return &Spec{
			id:    "message_length_limit",
			name:  "Message length limit",
			maker: newMessageLengthLimit,
		}
	})
}

var _ pkg.Module = &MessageLengthLimit{}

type MessageLengthLimit struct {
	mbase.Base
}

func newMessageLengthLimit(b mbase.Base) pkg.Module {
	return &MessageLengthLimit{
		Base: b,
	}
}

func (m MessageLengthLimit) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return nil

	/*
		user := event.User
		message := event.Message

		if user.GetName() == "gazatu2" {
			return nil
		}

		if user.GetName() == "supibot" {
			return nil
		}

		messageLength := len(message.GetText())
		if messageLength > 140 {
			if messageLength > 420 {
				return twitchactions.DoTimeout(event.User, 600*time.Second, "Your message is way too long")
			}

			return twitchactions.DoTimeout(event.User, 300*time.Second, "Your message is too long, shorten it")
		}

		return nil
	*/
}
