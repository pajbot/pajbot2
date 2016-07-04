package apirequest

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/common"
)

type bttvAPI struct {
	Emotes []map[string]interface{} `json:"emotes"`
}

func LoadBttvEmotes(channel string) ([]common.Emote, error) {
	var url string
	if channel == "global" {
		url = "https://api.betterttv.net/emotes"
	} else {
		url = fmt.Sprintf("https://api.betterttv.net/2/channels/%s", channel)
	}
	blob, err := HTTPRequest(url)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var data bttvAPI
	err = json.Unmarshal(blob, &data)
	if err != nil {
		log.Error(err)
	}
	if data.Emotes == nil {
		log.Error("no data")
	}
	if channel == "global" {
		return globalEmotes(data), nil
	}
	return channelEmotes(data), nil
}

func globalEmotes(data bttvAPI) []common.Emote {
	var emotes []common.Emote
	for _, e := range data.Emotes {
		name := e["regex"].(string)
		spl := strings.Split(e["url"].(string), "/emote/")[1]
		id := spl[:len(spl)-3] // remove /1x
		sizeX := e["width"].(float64)
		sizeY := e["height"].(float64)
		isGif := e["imageType"].(string) == "gif"
		emote := common.Emote{
			Name:  name,
			ID:    id,
			Type:  "bttv",
			SizeX: int(sizeX),
			SizeY: int(sizeY),
			IsGif: isGif,
			Count: 1,
		}
		emotes = append(emotes, emote)
	}
	return emotes
}

func channelEmotes(data bttvAPI) []common.Emote {
	var emotes []common.Emote
	for _, e := range data.Emotes {
		name := e["code"].(string)
		id := e["id"].(string)
		isGif := e["imageType"].(string) == "gif"
		emote := common.Emote{
			Name:  name,
			ID:    id,
			Type:  "bttv",
			SizeX: 28,
			SizeY: 28,
			IsGif: isGif,
			Count: 1,
		}
		emotes = append(emotes, emote)
	}
	return emotes
}
