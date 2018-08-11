package pkg

type ActionType interface {
	Do(Sender, Channel, User) error
	Priority() int
}

type Action interface {
	Do() error
	Set(ActionType)
}

var _ Action = &TwitchAction{}

type TwitchAction struct {
	Sender Sender

	Channel Channel

	User User

	action ActionType
}

type Timeout struct {
	Duration int
	Reason   string
}

func (a Timeout) Do(sender Sender, channel Channel, user User) error {
	sender.Timeout(channel, user, a.Duration, a.Reason)

	return nil
}

func (a Timeout) Priority() int {
	return 100 + a.Duration
}

type Ban struct {
	Reason string
}

func (a Ban) Do(sender Sender, channel Channel, user User) error {
	sender.Ban(channel, user, a.Reason)

	return nil
}

func (a Ban) Priority() int {
	return 0
}

func (a TwitchAction) Do() error {
	if a.action != nil {
		return a.action.Do(a.Sender, a.Channel, a.User)
	}

	return nil
}

func (a *TwitchAction) Set(action ActionType) {
	if a.action == nil {
		a.action = action
	} else {
		if a.action.Priority() > action.Priority() {
			a.action = action
		}
	}
}
