package apirequest

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	twitchTimeFormat = "2006-01-02T15:04:05Z"
	twitchAPIUrl     = "https://api.twitch.tv/kraken/%s"
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

type twitchStreamsChannelAPI struct {
	Stream  twitchAPIStream  `json:"stream"`
	Channel twitchAPIChannel `json:"channel"`
}

type twitch struct {
}

// TwitchAPI contains all methods relevant to the twitch api
var TwitchAPI = twitch{}

func (t *twitch) GetStream(channel string) (*Stream, error) {
	bs, err := HTTPRequest(fmt.Sprintf(twitchAPIUrl, "streams/"+channel))
	if err != nil {
		return nil, fmt.Errorf("HTTP Error %s", err)
	}
	var data twitchStreamsChannelAPI
	err = json.Unmarshal(bs, &data)
	if err != nil {
		return nil, fmt.Errorf("Invalid json: %s", err)
	}

	if data.Stream.ID == 0 {
		return nil, fmt.Errorf("Stream offline")
	}
	created, err := time.Parse(twitchTimeFormat, data.Stream.Created)
	if err != nil {
		return nil, fmt.Errorf("Invalid date format, probably malformed JSON")
	}
	stream := &Stream{}
	stream.ID = fmt.Sprintf("%.f", data.Stream.ID)
	stream.Online = !data.Stream.IsPlaylist
	stream.Created = created
	stream.Game = data.Stream.Game
	stream.Viewers = int(data.Stream.Viewers)
	stream.IsPlaylist = data.Stream.IsPlaylist
	return stream, nil
}
