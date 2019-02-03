package modules

import (
	"github.com/pajlada/pajbot2/pkg"
)

type badCharacterFilter struct {
	botChannel pkg.BotChannel

	badCharacters []rune
}

func newBadCharacterFilter() pkg.Module {
	return &badCharacterFilter{}
}

var badCharacterSpec = moduleSpec{
	id:    "bad_character_filter",
	name:  "Bad character filter",
	maker: newBadCharacterFilter,
}

func (m *badCharacterFilter) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	m.badCharacters = append(m.badCharacters, '\x01')

	return nil
}

func (m *badCharacterFilter) Disable() error {
	return nil
}

func (m *badCharacterFilter) Spec() pkg.ModuleSpec {
	return &badCharacterSpec
}

func (m *badCharacterFilter) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *badCharacterFilter) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *badCharacterFilter) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	for _, r := range message.GetText() {
		for _, badCharacter := range m.badCharacters {
			if r == badCharacter {
				action.Set(pkg.Timeout{
					Duration: 300, Reason: "Your message contains a banned character",
				})
				return nil
			}
		}
	}

	return nil
}
