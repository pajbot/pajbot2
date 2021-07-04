package apirequest

import (
	"github.com/pajbot/pajbot2/pkg/common/config"
)

func InitTwitch(cfg *config.Config) (err error) {
	err = initWrapper(&cfg.Auth.Twitch)

	return
}
