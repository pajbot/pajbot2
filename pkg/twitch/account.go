package twitch

type SimpleAccount struct {
	id   string
	name string
}

func (a *SimpleAccount) ID() string {
	return a.id
}

func (a *SimpleAccount) Name() string {
	return a.name
}

type TwitchAccount struct {
	UserID   string
	UserName string
}

func (a *TwitchAccount) ID() string {
	return a.UserID
}

func (a *TwitchAccount) Name() string {
	return a.UserName
}
