package pkg

type Message interface {
	GetText() string
	SetText(string)

	GetTwitchReader() EmoteReader

	GetBTTVReader() EmoteReader
	AddBTTVEmote(Emote)
}
