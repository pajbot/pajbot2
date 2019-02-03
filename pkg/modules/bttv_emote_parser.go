package modules

import (
	"strings"
	"unicode"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/common"
	"github.com/pajlada/pajbot2/pkg/emotes"
)

type bttvEmoteParser struct {
	botChannel pkg.BotChannel

	globalEmotes *map[string]common.Emote
}

var bttvEmoteParserSpec = &moduleSpec{
	id:   "bttv_emote_parser",
	name: "BTTV emote parser",

	maker: newbttvEmoteParser,

	enabledByDefault: true,

	priority: -50000,
}

func newbttvEmoteParser() pkg.Module {
	return &bttvEmoteParser{
		globalEmotes: &emotes.GlobalEmotes.Bttv,
	}
}

func (m *bttvEmoteParser) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel
	return nil
}

func (m *bttvEmoteParser) Disable() error {
	return nil
}

func (m *bttvEmoteParser) Spec() pkg.ModuleSpec {
	return bttvEmoteParserSpec
}

func (m *bttvEmoteParser) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m *bttvEmoteParser) OnWhisper(bot pkg.BotChannel, user pkg.User, message pkg.Message) error {
	return nil
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
