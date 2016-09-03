package common

// Emote xD
type Emote struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	// Possible types: bttv, twitch
	Type string `json:"type"`

	// Size in pixels
	SizeX int `json:"size_x"`
	SizeY int `json:"size_y"`

	IsGif bool `json:"is_gif"`
	Count int  `json:"count"`
}
