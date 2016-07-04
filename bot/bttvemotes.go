package bot

import (
	_ "log" // go-imports pajaSWA
	"strings"

	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/common"
)

type bttvAPI struct {
	Emotes []map[string]interface{} `json:"emotes"`
}

// LoadBttvEmotes should load emotes from redis, but this should do for now
func (bot *Bot) LoadBttvEmotes() {
	channelEmotes, err := apirequest.LoadBttvEmotes(bot.Channel.Name)
	if err != nil {
		log.Error(err)
		return
	}
	for _, emote := range channelEmotes {
		bot.Channel.BttvEmotes[emote.Name] = emote
	}
	globalEmotes, err := apirequest.LoadBttvEmotes("global")
	if err != nil {
		log.Error(err)
		return
	}
	for _, emote := range globalEmotes {
		bot.Channel.BttvEmotes[emote.Name] = emote
	}
}

// regex would probably be better but im a regex noob ¯\_(ツ)_/¯
func (bot *Bot) parseBttvEmotes(msg *common.Msg) {
	m := strings.Split(msg.Text, " ")
	for _, word := range m {
		if emote, ok := bot.Channel.BttvEmotes[word]; ok {
			msg.Emotes = append(msg.Emotes, emote)
			log.Debug(emote)
		}
	}
}
