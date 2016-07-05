package apirequest

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	TwitchTimeFormat = "2006-01-02T15:04:05Z"
	TwitchAPIUrl     = "https://api.twitch.tv/kraken/%s"
)

type twitchAPIStream struct {
	ID         float64 `json:"_id"`
	Game       string  `json:"game"`
	Viewers    float64 `json:"viewers"`
	Created    string  `json:"created_at"`
	IsPlaylist bool    `json:"is_playlist"`
}

type twitchAPIChannel struct {
	ID        float64 `json:"_id"`
	Status    string  `json:"status"`
	Game      string  `json:"game"`
	UpdatedAt string  `json:"updated_at"`
	Views     float64 `json:"views"`
	Followers float64 `json:"followers"`
	Partner   bool    `json:"partner"`
}

type twitchAPI struct {
	Stream  twitchAPIStream  `json:"stream"`
	Channel twitchAPIChannel `json:"channel"`
}

func GetStream(channel string) (Stream, error) {
	var stream Stream
	bs, err := HTTPRequest(fmt.Sprintf(TwitchAPIUrl, "streams/"+channel))
	if err != nil {
		return stream, err
	}
	var data twitchAPI
	err = json.Unmarshal(bs, &data)
	if err != nil {
		log.Error(err)
		return stream, err
	}
	log.Debug(data)

	if data.Stream.ID == 0 {
		// stream offline
		return stream, nil
	}
	created, err := time.Parse(TwitchTimeFormat, data.Stream.Created)
	if err != nil {
		log.Error(err)
	}
	stream.ID = fmt.Sprintf("%.f", data.Stream.ID)
	stream.Online = !data.Stream.IsPlaylist
	stream.Created = created
	stream.Game = data.Stream.Game
	stream.Viewers = int(data.Stream.Viewers)
	stream.IsPlaylist = data.Stream.IsPlaylist
	return stream, nil
}
