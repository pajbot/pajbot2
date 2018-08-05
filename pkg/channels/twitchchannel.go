package channels

type TwitchChannel struct {
	Channel string
	ID      string
}

func (c TwitchChannel) GetChannel() string {
	return c.Channel
}

func (c TwitchChannel) GetID() string {
	return c.ID
}
