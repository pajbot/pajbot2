package bad_character_filter

import (
	"time"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

func init() {
	modules.Register("bad_character_filter", func() pkg.ModuleSpec {
		return modules.NewSpec("bad_character_filter", "Bad character filter", false, newBadCharacterFilter)
	})
}

type badCharacterFilter struct {
	mbase.Base

	badCharacters []rune
}

func newBadCharacterFilter(b *mbase.Base) pkg.Module {
	return &badCharacterFilter{
		Base: *b,

		badCharacters: []rune{'\x01'},
	}
}

func (m *badCharacterFilter) OnMessage(event pkg.MessageEvent) pkg.Actions {
	message := event.Message

	for _, r := range message.GetText() {
		for _, badCharacter := range m.badCharacters {
			if r == badCharacter {
				actions := &twitchactions.Actions{}
				actions.Timeout(event.User, 300*time.Second).SetReason("Your message contains a banned character")
				return actions
			}
		}
	}

	return nil
}
