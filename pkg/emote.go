package pkg

type Emote interface {
	// i.e. "8vn23893nvcakuj23" for a bttv emote, or "85489481" for a twitch emote
	GetID() string

	// i.e. "NaM" or "forsenE"
	GetName() string

	// "twitch" or "bttv"
	GetType() string

	GetCount() int
}

type EmoteReader interface {
	Next() bool
	Get() Emote
}
