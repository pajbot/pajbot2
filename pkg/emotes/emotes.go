package emotes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pajbot/pajbot2/pkg/apirequest"
	"github.com/pajbot/pajbot2/pkg/common"
	"github.com/pajlada/gobttv"
	"github.com/pajlada/goffz"
)

// GlobalEmotes contains BTTV & FFZ Emotes
var GlobalEmotes common.ExtensionEmotes

func loadGlobalBttvEmotes() error {
	emotes, err := apirequest.BTTV.GetEmotes()
	if err != nil {
		return err
	}

	GlobalEmotes.Bttv = make(map[string]common.Emote)
	GlobalEmotes.BttvLastUpdate = time.Now()

	for _, emote := range emotes {
		GlobalEmotes.Bttv[emote.Code] = ParseBTTVGlobalEmote(emote)
	}

	return nil
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
			fmt.Printf("Error fetching Global FFZ Emotes")
			fmt.Printf("Status code: %d", statusCode)
			fmt.Printf("Status message: %s", statusMessage)
			fmt.Printf("Error message: %s", errorMessage)
		}, func(err error) {
			fmt.Printf("Internal error: %s", err)
		})
}

// LoadGlobalEmotes xD
func LoadGlobalEmotes() {
	loadGlobalBttvEmotes()
	loadGlobalFrankerFaceZEmotes()
}

// ParseBTTVGlobalEmote parses a BTTV emote into a common.Emote
func ParseBTTVGlobalEmote(emote gobttv.Emote) common.Emote {
	return common.Emote{
		Name: emote.Code,
		ID:   emote.ID,
		Type: "bttv",
		// SizeX: emote.Width,
		// SizeY: emote.Height,
		SizeX: 28, // This data is no longer provided by the BTTV api, so this is inaccurate
		SizeY: 28, // This data is no longer provided by the BTTV api, so this is inaccurate
		Count: 1,
		IsGif: emote.ImageType == "gif",
	}
}

// ParseBTTVChannelEmote parses a BTTV emote into a common.Emote
func ParseBTTVChannelEmote(emote gobttv.Emote) common.Emote {
	return common.Emote{
		Name:  emote.Code,
		ID:    emote.ID,
		Type:  "bttv",
		SizeX: 28,
		SizeY: 28,
		Count: 1,
		IsGif: emote.ImageType == "gif",
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
