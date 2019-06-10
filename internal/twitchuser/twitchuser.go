package twitchuser

type TwitchUser struct {
	id   string
	name string
}

func New(id, name string) *TwitchUser {
	return &TwitchUser{
		id:   id,
		name: name,
	}
}

func (u *TwitchUser) ID() string {
	return u.id
}

func (u *TwitchUser) Name() string {
	return u.name
}
