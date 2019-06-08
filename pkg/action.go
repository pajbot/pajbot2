package pkg

import "fmt"

type ActionType int

const (
	ActionTypeTimeout ActionType = iota
	ActionTypeBan
)

type ActionPerformer interface {
	Do(Sender, Channel, User) error
	Priority() int
	Type() ActionType
}

type Action interface {
	Do() error
	Set(ActionPerformer)
	Get() ActionPerformer

	SetReason(string)
	Reason() string

	NotifyModerator() User
	SetNotifyModerator(User)
}

var _ Action = &TwitchAction{}

type TwitchAction struct {
	Sender Sender

	Channel Channel

	User User

	Soft bool

	action ActionPerformer

	reason string

	notifyModerator User
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

func (a Timeout) Type() ActionType {
	return ActionTypeTimeout
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

func (a Ban) Type() ActionType {
	return ActionTypeBan
}

func (a TwitchAction) Do() error {
	if a.Soft {
		return nil
	}

	if a.action != nil {
		if a.NotifyModerator() != nil {
			a.Sender.Whisper(a.NotifyModerator(), fmt.Sprintf("%s triggered bad banphrase in %s", a.User.GetName(), a.Channel.GetName()))
		}
		return a.action.Do(a.Sender, a.Channel, a.User)
	}

	return nil
}

func (a *TwitchAction) Set(action ActionPerformer) {
	if a.action == nil {
		a.action = action
	} else {
		if action.Priority() > a.action.Priority() {
			a.action = action
		}
	}
}

func (a *TwitchAction) Get() ActionPerformer {
	return a.action
}

func (a *TwitchAction) SetReason(reason string) {
	a.reason = reason
}

func (a *TwitchAction) Reason() string {
	return a.reason
}

func (a TwitchAction) NotifyModerator() User {
	return a.notifyModerator
}

func (a *TwitchAction) SetNotifyModerator(user User) {
	a.notifyModerator = user
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

type A interface {
	// timeouts
	// Timeout(username, duration)

	// messages
	// Say(message)

	// DoOnSuccess(func())

	// Do(func())
}
