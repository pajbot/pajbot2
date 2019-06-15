package twitchactions

import "github.com/pajbot/pajbot2/pkg"

var _ pkg.MessageAction = &Message{}

type Message struct {
	content string

	action bool
}

func (m *Message) SetAction(v bool) {
	m.action = v
}

func (m Message) Evaluate() string {
	if m.action {
		return ".me " + m.content
	}

	return m.content
}
