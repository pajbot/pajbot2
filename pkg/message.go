package pkg

type Message interface {
	GetText() string

	GetTwitchReader() EmoteReader

	GetBTTVReader() EmoteReader
	AddBTTVEmote(Emote)
}
