package points

import (
	"errors"
	"strconv"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// GivePoints xD
func GivePoints(b *bot.Bot, user *common.User, args []string) error {
	if len(args) < 2 {
		return errors.New("not enough args xD")
	}
	pts, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}
	if !b.Redis.IsValidUser(b.Channel.Name, args[0]) {
		return errors.New("invalid user")
	}
	b.Redis.IncrPoints(b.Channel.Name, args[0], int(pts))
	return nil
}
