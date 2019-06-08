package modules

import (
	"strings"
	"unicode"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common"
	"github.com/pajbot/pajbot2/pkg/emotes"
)

func init() {
	Register("bttv_emote_parser", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:   "bttv_emote_parser",
			name: "BTTV emote parser",

			enabledByDefault: true,

			priority: -50000,

			maker: newbttvEmoteParser,
		}
	})
}

type bttvEmoteParser struct {
	base

	globalEmotes *map[string]common.Emote
}

func newbttvEmoteParser(b base) pkg.Module {
	return &bttvEmoteParser{
		base: b,

		globalEmotes: &emotes.GlobalEmotes.Bttv,
	}
}

func (m *bttvEmoteParser) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.FieldsFunc(message.GetText(), func(r rune) bool {
		// TODO(pajlada): This needs better testing
		return r > 0xFF || unicode.IsSpace(r) || r == '!' || r == '.' || r == '$' || r == '^' || r == '#' || r == '*' || r == '@' || r == ')' || r == '%' || r == '&' || r > 0x7a || r < 0x30 || (r > 0x39 && r < 0x41) || (r > 0x5a && r < 0x61)
	})
	emoteCount := make(map[string]*common.Emote)
	for _, word := range parts {
		if emote, ok := emoteCount[word]; ok {
			emote.Count++
		} else if emote, ok := (*m.globalEmotes)[word]; ok {
			emoteCount[word] = &emote
		}
	}

	for _, emote := range emoteCount {
		message.AddBTTVEmote(emote)
	}

	return nil
}
