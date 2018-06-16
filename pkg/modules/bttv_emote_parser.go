package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/emotes"
	"github.com/pajlada/pajbot2/pkg"
)

type BTTVEmoteParser struct {
	server *server
}

func NewBTTVEmoteParser() *BTTVEmoteParser {
	return &BTTVEmoteParser{
		server: &_server,
	}
}

func (m BTTVEmoteParser) Name() string {
	return "BTTVEmoteParser"
}

func (m BTTVEmoteParser) Register() error {
	return nil
}

func (m BTTVEmoteParser) OnWhisper(user pkg.User, message pkg.Message) error {
	return nil
}

func (m BTTVEmoteParser) OnMessage(channel pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	cleanText := strings.Map(
		func(r rune) rune {
			if r <= 0xFF {
				return r
			}

			return -1
		}, message.GetText())

	// parts := strings.Split(message.GetText(), " ")
	// parts := strings.Fields(message.GetText())
	parts := strings.Fields(cleanText)
	emoteCount := make(map[string]*common.Emote)
	for _, word := range parts {
		if emote, ok := emoteCount[word]; ok {
			emote.Count++
		} else if emote, ok := emotes.GlobalEmotes.Bttv[word]; ok {
			emoteCount[word] = &emote
		}
	}

	for _, emote := range emoteCount {
		message.AddBTTVEmote(emote)
	}

	return nil
}
