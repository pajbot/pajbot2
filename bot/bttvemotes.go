package bot

import (
	"encoding/json"
	"io/ioutil"
	_ "log" // go-imports pajaSWA
	"net/http"
	"strings"

	"github.com/pajlada/pajbot2/common"
)

type bttvAPI struct {
	Emotes []map[string]interface{} `json:"emotes"`
}

// LoadBttvEmotes should load emotes from redis, but this should do for now
func (bot *Bot) LoadBttvEmotes() {
	req, err := http.Get("https://api.betterttv.net/emotes")
	if err != nil {
		log.Error(err)
	}
	blob, _ := ioutil.ReadAll(req.Body)
	var data bttvAPI
	err = json.Unmarshal(blob, &data)
	if err != nil {
		log.Error(err)
	}
	if data.Emotes == nil {
		log.Error("no data")
	}
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
		bot.Channel.BttvEmotes[name] = emote
		//log.Debug(emote)
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
