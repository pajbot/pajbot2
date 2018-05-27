package channels

type TwitchChannel struct {
	Channel string
}

func (c TwitchChannel) GetChannel() string {
	return c.Channel
}
