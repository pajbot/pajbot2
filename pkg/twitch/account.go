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
