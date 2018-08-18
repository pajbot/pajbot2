package modules

import (
	"strings"
	"unicode"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/common"
)

type BTTVEmoteParser struct {
	server *server

	globalEmotes *map[string]common.Emote
}

func NewBTTVEmoteParser(globalEmotes *map[string]common.Emote) *BTTVEmoteParser {
	return &BTTVEmoteParser{
		server:       &_server,
		globalEmotes: globalEmotes,
	}
}

func (m BTTVEmoteParser) Name() string {
	return "BTTVEmoteParser"
}

func (m BTTVEmoteParser) Register() error {
	return nil
}

func (m BTTVEmoteParser) OnWhisper(bot pkg.Sender, user pkg.User, message pkg.Message) error {
	return nil
}

func (m BTTVEmoteParser) OnMessage(bot pkg.Sender, channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
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
