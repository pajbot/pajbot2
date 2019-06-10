package twitchactions

import (
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.MuteAction = &Ban{}

type Ban struct {
	mute

	duration time.Duration
}

func (b Ban) Type() pkg.MuteType {
	return pkg.MuteTypePermanent
}

func (b Ban) Duration() time.Duration {
	return time.Duration(0)
}
