package pkg

type PubSubAuthorization struct {
	Nonce        string
	TwitchUserID string
	admin        bool
}

func PubSubAdminAuth() *PubSubAuthorization {
	return &PubSubAuthorization{
		admin: true,
	}
}

func (p PubSubAuthorization) Admin() bool {
	return p.admin
}

type PubSubBan struct {
	Channel string
	Target  string
	Reason  string
}

type PubSubUntimeout struct {
	Channel string
	Target  string
}
