package points

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

func Roulette(b *bot.Bot, user *common.User, args []string) error {
	if user.Points < 1 {
		return fmt.Errorf("you dont have enough points to roulette %s ;p", user.Name)
	}
	if len(args) < 1 {
		return fmt.Errorf("usage: !roul 123")
	}
	if args[0] == "all" || args[0] == "allin" {
		runRoulette(b, user, user.Points)
		return nil
	}
	_bet, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("usage: !roul 123")
	}
	bet := int(_bet)
	if bet > user.Points {
		return fmt.Errorf("%s, you dont have that many points :p", user.Name)
	}
	if bet < 1 {
		return fmt.Errorf("%s, you cant roulette 0 points pajaSWA", user.Name)
	}
	runRoulette(b, user, bet)
	return nil
}

func runRoulette(b *bot.Bot, user *common.User, points int) {
	won := rand.Float32() >= 0.5
	if won {
		b.Redis.IncrPoints(b.Channel.Name, user.Name, points)
		b.Sayf("%s won %d points in roulette and now has %d points pajaDank",
			user.Name, points, user.Points+points)
	} else {
		b.Redis.IncrPoints(b.Channel.Name, user.Name, -points)
		b.Sayf("%s lost %d points in roulette and now has %d points LUL",
			user.Name, points, user.Points-points)
	}

}
