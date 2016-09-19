package bot

import (
	_ "log" // go-imports pajaSWA
	"strconv"
	"strings"
	"time"

	"github.com/pajlada/gobttv"
	"github.com/pajlada/goffz"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/common"
)

// GlobalEmotes contains BTTV & FFZ Emotes
var GlobalEmotes common.ExtensionEmotes

func loadGlobalBttvEmotes() {
	apirequest.BTTV.GetEmotes(
		func(emotesResponse gobttv.EmotesResponse) {
			GlobalEmotes.Bttv = make(map[string]common.Emote)
			GlobalEmotes.BttvLastUpdate = time.Now()

			for _, emote := range emotesResponse.Emotes {
				GlobalEmotes.Bttv[emote.Regex] = ParseBTTVGlobalEmote(emote)
			}
		},
		func(statusCode int, statusMessage, errorMessage string) {
			log.Errorf("Error fetching Global BTTV Emotes")
			log.Errorf("Status code: %d", statusCode)
			log.Errorf("Status message: %s", statusMessage)
			log.Errorf("Error message: %s", errorMessage)
		}, func(err error) {
			log.Errorf("Internal error: %s", err)
		})
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
	apirequest.BTTV.GetChannel(bot.Channel.Name,
		func(channel gobttv.ChannelResponse) {
			bot.Channel.Emotes.Bttv = make(map[string]common.Emote)
			bot.Channel.Emotes.BttvLastUpdate = time.Now()

			for _, emote := range channel.Emotes {
				bot.Channel.Emotes.Bttv[emote.Code] = ParseBTTVChannelEmote(emote)
			}
		},
		func(statusCode int, statusMessage, errorMessage string) {
			// We ignore 404 errors, it just means he doesn't have a BTTV account
			if statusCode != 404 {
				log.Errorf("Error fetching Channel BTTV Emotes (%s)", bot.Channel.Name)
				log.Errorf("Status code: %d", statusCode)
				log.Errorf("Status message: %s", statusMessage)
				log.Errorf("Error message: %s", errorMessage)
			}
		}, func(err error) {
			log.Errorf("Internal error: %s", err)
		})
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
	maxScale := 1
	var keyScale int
	var err error
	for k := range emote.Urls {
		keyScale, err = strconv.Atoi(k)
		if err == nil && keyScale > maxScale {
			maxScale = keyScale
		}
	}
	return common.Emote{
		Name:     emote.Name,
		ID:       strconv.Itoa(emote.ID),
		Type:     "ffz",
		SizeX:    emote.Width,
		SizeY:    emote.Height,
		Count:    1,
		MaxScale: maxScale,
	}
}

// ParseBTTVGlobalEmote parses a BTTV emote into a common.Emote
func ParseBTTVGlobalEmote(emote gobttv.GlobalEmoteData) common.Emote {
	spl := strings.Split(emote.URL, "/emote/")[1]
	id := spl[:len(spl)-3] // remove /1x
	isGif := emote.ImageType == "gif"
	return common.Emote{
		Name:  emote.Regex,
		ID:    id,
		Type:  "bttv",
		SizeX: emote.Width,
		SizeY: emote.Height,
		Count: 1,
		IsGif: isGif,
	}
}

// ParseBTTVChannelEmote parses a BTTV emote into a common.Emote
func ParseBTTVChannelEmote(emote gobttv.ChannelEmoteData) common.Emote {
	isGif := emote.ImageType == "gif"
	return common.Emote{
		Name:  emote.Code,
		ID:    emote.ID,
		Type:  "bttv",
		SizeX: 28,
		SizeY: 28,
		Count: 1,
		IsGif: isGif,
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
