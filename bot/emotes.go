package bot

import (
	_ "log" // go-imports pajaSWA
	"strconv"
	"strings"
	"time"

	"github.com/pajlada/goffz"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/common"
)

// GlobalEmotes contains BTTV & FFZ Emotes
var GlobalEmotes common.ExtensionEmotes

func loadGlobalBttvEmotes() {
	globalEmotes, err := apirequest.BTTVAPI.LoadEmotes("global")
	if err != nil {
		log.Error(err)
		return
	}
	GlobalEmotes.Bttv = make(map[string]common.Emote)
	for _, emote := range globalEmotes {
		GlobalEmotes.Bttv[emote.Name] = emote
	}
	GlobalEmotes.BttvLastUpdate = time.Now()
}

func loadGlobalFrankerFaceZEmotes() {
	apirequest.FFZ.GetSet("global",
		func(rSet goffz.SetResponse) {
			GlobalEmotes.FrankerFaceZ = make(map[string]common.Emote)
			GlobalEmotes.FrankerFaceZLastUpdate = time.Now()

			for _, set := range rSet.Sets {
				for _, emote := range set.Emoticons {
					GlobalEmotes.FrankerFaceZ[emote.Name] = ParseFrankerFaceZEmote(emote)
				}
			}
		},
		func(statusCode int, statusMessage, errorMessage string) {
			log.Errorf("Error fetching Global FFZ Emotes")
			log.Errorf("Status code: %d", statusCode)
			log.Errorf("Status message: %s", statusMessage)
			log.Errorf("Error message: %s", errorMessage)
		}, func(err error) {
			log.Errorf("Internal error: %s", err)
		})
}

// LoadGlobalEmotes xD
func LoadGlobalEmotes() {
	go loadGlobalBttvEmotes()
	go loadGlobalFrankerFaceZEmotes()
}

// LoadBttvEmotes should load emotes from redis, but this should do for now
func (bot *Bot) LoadBttvEmotes() {
	channelEmotes, err := apirequest.BTTVAPI.LoadEmotes(bot.Channel.Name)
	if err != nil {
		log.Error(err)
		return
	}
	bot.Channel.Emotes.BttvLastUpdate = time.Now()

	bot.Channel.Emotes.Bttv = make(map[string]common.Emote)
	for _, emote := range channelEmotes {
		bot.Channel.Emotes.Bttv[emote.Name] = emote
	}
}

// regex would probably be better but im a regex noob ¯\_(ツ)_/¯
func (bot *Bot) parseEmotes(msg *common.Msg) {
	m := strings.Split(msg.Text, " ")
	emoteCount := make(map[string]*common.Emote)
	for _, word := range m {
		if emote, ok := emoteCount[word]; ok {
			emote.Count++
		} else if emote, ok := bot.Channel.Emotes.Bttv[word]; ok {
			emoteCount[word] = &emote
		} else if emote, ok := GlobalEmotes.Bttv[word]; ok {
			emoteCount[word] = &emote
		} else if emote, ok := bot.Channel.Emotes.FrankerFaceZ[word]; ok {
			emoteCount[word] = &emote
		} else if emote, ok := GlobalEmotes.FrankerFaceZ[word]; ok {
			emoteCount[word] = &emote
		}
	}

	for _, emote := range emoteCount {
		msg.Emotes = append(msg.Emotes, *emote)
	}
}

// ParseFrankerFaceZEmote parses a FFZ emote into a common.Emote
func ParseFrankerFaceZEmote(emote goffz.EmoteData) common.Emote {
	return common.Emote{
		Name:  emote.Name,
		ID:    strconv.Itoa(emote.ID),
		Type:  "ffz",
		SizeX: emote.Width,
		SizeY: emote.Height,
		Count: 1,
	}
}

// LoadFFZEmotes should load emotes from redis, but this should do for now
func (bot *Bot) LoadFFZEmotes() {
	apirequest.FFZ.GetRoom(bot.Channel.Name,
		func(room goffz.RoomResponse) {
			bot.Channel.Emotes.FrankerFaceZ = make(map[string]common.Emote)
			bot.Channel.Emotes.FrankerFaceZLastUpdate = time.Now()

			for _, set := range room.Sets {
				for _, emote := range set.Emoticons {
					bot.Channel.Emotes.FrankerFaceZ[emote.Name] = ParseFrankerFaceZEmote(emote)
				}
			}
		},
		func(statusCode int, statusMessage, errorMessage string) {
			// We ignore 404 errors, it just means he doesn't have a FFZ account
			if statusCode != 404 {
				log.Errorf("Error fetching Channel FFZ Emotes (%s)", bot.Channel.Name)
				log.Errorf("Status code: %d", statusCode)
				log.Errorf("Status message: %s", statusMessage)
				log.Errorf("Error message: %s", errorMessage)
			}
		}, func(err error) {
			log.Errorf("Internal error: %s", err)
		})
}
