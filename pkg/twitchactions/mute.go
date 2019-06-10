package twitchactions

import "github.com/pajbot/pajbot2/pkg"

type mute struct {
	user   pkg.User
	reason string
}

func (m *mute) User() pkg.User {
	return m.user
}

func (m *mute) SetReason(reason string) {
	m.reason = reason
}

func (m mute) Reason() string {
	return m.reason
}
