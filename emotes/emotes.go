package emotes

import (
	"log"
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
	log.Println("XD BTTV??")
	apirequest.BTTV.GetEmotes(
		func(emotesResponse gobttv.EmotesResponse) {
			log.Println("XD BTTV")
			GlobalEmotes.Bttv = make(map[string]common.Emote)
			GlobalEmotes.BttvLastUpdate = time.Now()

			for _, emote := range emotesResponse.Emotes {
				GlobalEmotes.Bttv[emote.Regex] = ParseBTTVGlobalEmote(emote)
			}
		},
		func(statusCode int, statusMessage, errorMessage string) {
			log.Printf("Error fetching Global BTTV Emotes")
			log.Printf("Status code: %d", statusCode)
			log.Printf("Status message: %s", statusMessage)
			log.Printf("Error message: %s", errorMessage)
		}, func(err error) {
			log.Printf("Internal error: %s", err)
		})
}

func loadGlobalFrankerFaceZEmotes() {
	log.Println("Loading FFZ Emotes...")
	apirequest.FFZ.GetSet("global",
		func(rSet goffz.SetResponse) {
			log.Println("Done loading FFZ emotes!")
			GlobalEmotes.FrankerFaceZ = make(map[string]common.Emote)
			GlobalEmotes.FrankerFaceZLastUpdate = time.Now()

			for _, set := range rSet.Sets {
				for _, emote := range set.Emoticons {
					GlobalEmotes.FrankerFaceZ[emote.Name] = ParseFrankerFaceZEmote(emote)
				}
			}
		},
		func(statusCode int, statusMessage, errorMessage string) {
			log.Printf("Error fetching Global FFZ Emotes")
			log.Printf("Status code: %d", statusCode)
			log.Printf("Status message: %s", statusMessage)
			log.Printf("Error message: %s", errorMessage)
		}, func(err error) {
			log.Printf("Internal error: %s", err)
		})
}

// LoadGlobalEmotes xD
func LoadGlobalEmotes() {
	loadGlobalBttvEmotes()
	loadGlobalFrankerFaceZEmotes()
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
