package modules

import (
	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("bad_character_filter", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "bad_character_filter",
			name:  "Bad character filter",
			maker: newBadCharacterFilter,
		}
	})
}

type badCharacterFilter struct {
	base

	badCharacters []rune
}

func newBadCharacterFilter(b base) pkg.Module {
	return &badCharacterFilter{
		base: b,

		badCharacters: []rune{'\x01'},
	}
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
