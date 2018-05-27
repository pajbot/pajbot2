package modules

import (
	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/pkg"
)

type BadCharacterFilter struct {
	server *server

	badCharacters []rune

	Sender pkg.Sender
}

func NewBadCharacterFilter(sender pkg.Sender) *BadCharacterFilter {
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

func (m BadCharacterFilter) OnWhisper(source pkg.User, message twitch.Message) error {
	return nil
}

func (m BadCharacterFilter) OnMessage(source pkg.Channel, user pkg.User, message twitch.Message) error {
	for _, r := range message.Text {
		for _, badCharacter := range m.badCharacters {
			if r == badCharacter {
				m.Sender.Timeout(source, user, 300, "Message contains a banned character")
				return nil
			}
		}
	}

	return nil
}
