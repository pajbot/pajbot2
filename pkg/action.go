package pkg

type Timeout struct {
	duration int
	reason   string
}

type Action interface {
	Do() error
	SetTimeout(int, string)
}

var _ Action = &TwitchAction{}

type TwitchAction struct {
	Sender Sender

	Channel Channel

	User User

	timeout Timeout
}

func (a TwitchAction) Do() error {
	if a.timeout.duration > 0 {
		a.Sender.Timeout(a.Channel, a.User, a.timeout.duration, a.timeout.reason)
	}

	return nil
}

func (a *TwitchAction) SetTimeout(duration int, reason string) {
	if duration > a.timeout.duration {
		a.timeout.duration = duration
		a.timeout.reason = reason
	}
}
