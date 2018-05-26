package modules

import (
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/pkg"
)

type BadCharacterFilter struct {
	server *server

	badCharacters []rune

	Sender pkg.Channel
}

func NewBadCharacterFilter(sender pkg.Channel) *BadCharacterFilter {
	return &BadCharacterFilter{
		server: &_server,
		Sender: sender,
	}
}

func (m *BadCharacterFilter) Register() error {
	m.badCharacters = append(m.badCharacters, '\x01')

	return nil
}

func (m BadCharacterFilter) Name() string {
	return "BadCharacterFilter"
}

func (m BadCharacterFilter) OnMessage(channel string, user twitch.User, message twitch.Message) error {
	for _, r := range message.Text {
		for _, badCharacter := range m.badCharacters {
			if r == badCharacter {
				m.Sender.Timeout(channel, user, 300, "Message contains a banned character")
				return nil
			}
		}
	}

	return nil
}
