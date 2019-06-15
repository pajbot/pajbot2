package twitchactions

import (
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.MuteAction = &Timeout{}

type Timeout struct {
	mute

	duration time.Duration
}

func (t Timeout) Type() pkg.MuteType {
	return pkg.MuteTypeTemporary
}

func (t Timeout) Duration() time.Duration {
	return t.duration
}
