package common

import "time"

// Emote xD
type Emote struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	// Possible types: bttv, twitch, ffz
	Type string `json:"type"`

	// Size in pixels
	SizeX int `json:"size_x"`
	SizeY int `json:"size_y"`

	IsGif bool `json:"is_gif"`
	Count int  `json:"count"`

	MaxScale int `json:"max_scale"`
}

// ExtensionEmotes is an object which contains emotes that are shared between all channels
type ExtensionEmotes struct {
	// Global BTTV Emotes
	Bttv           map[string]Emote
	BttvLastUpdate time.Time

	// Global FrankerFaceZ Emotes
	FrankerFaceZ           map[string]Emote
	FrankerFaceZLastUpdate time.Time
}

// EmoteByName implements sort.Interface by emote name
type EmoteByName []Emote

// Len implements sort.Interface
func (a EmoteByName) Len() int {
	return len(a)
}

// Swap implements sort.Interface
func (a EmoteByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less implements sort.Interface
func (a EmoteByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func (e Emote) GetID() string {
	return e.ID
}

func (e Emote) GetName() string {
	return e.Name
}

func (e Emote) GetType() string {
	return e.Type
}

func (e Emote) GetCount() int {
	return e.Count
}
