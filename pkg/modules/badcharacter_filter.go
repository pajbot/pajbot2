package modules

import (
	"github.com/pajlada/pajbot2/pkg"
)

type BadCharacterFilter struct {
	server *server

	badCharacters []rune
}

func NewBadCharacterFilter() *BadCharacterFilter {
	return &BadCharacterFilter{
		server: &_server,
	}
}

func (m *BadCharacterFilter) Register() error {
	m.badCharacters = append(m.badCharacters, '\x01')

	return nil
}

func (m BadCharacterFilter) Name() string {
	return "BadCharacterFilter"
}

func (m BadCharacterFilter) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m BadCharacterFilter) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	for _, r := range message.GetText() {
		for _, badCharacter := range m.badCharacters {
			if r == badCharacter {
				action.SetTimeout(300, "Your message contains a banned character")
				return nil
			}
		}
	}

	return nil
}
