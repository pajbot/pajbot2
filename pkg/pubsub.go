package pkg

type PubSubBan struct {
	Channel string
	Target  string
	Reason  string
}

type PubSubUntimeout struct {
	Channel string
	Target  string
}
