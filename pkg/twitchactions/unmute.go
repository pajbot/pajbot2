package twitchactions

import "github.com/pajbot/pajbot2/pkg"

type unmute struct {
	user     pkg.User
	muteType pkg.MuteType
}

func (m *unmute) User() pkg.User {
	return m.user
}

func (m *unmute) Type() pkg.MuteType {
	return m.muteType
}
