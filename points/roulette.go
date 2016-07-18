package points

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// Roulette xD
type Roulette struct {
	WinMessage  string
	LoseMessage string
}

// Run roulette
func (r *Roulette) Run(b *bot.Bot, msg *common.Msg, args []string) error {
	user := &msg.User
	if user.Points == 0 {
		return fmt.Errorf("you dont have enough points to roulette %s ;p", user.Name)
	}
	if len(args) < 1 {
		return fmt.Errorf("usage: !roul 123")
	}
	if args[0] == "all" || args[0] == "allin" {
		r.runRoulette(b, msg, user.Points)
		return nil
	}
	_bet, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("usage: !roul 123")
	}
	bet := int(_bet)
	if bet == 0 {
		return fmt.Errorf("%s, you cant roulette 0 points pajaSWA", user.Name)
	}
	if (bet > 0 && bet > user.Points) || (bet < 0 && bet < user.Points) {
		return fmt.Errorf("%s, you dont have that many points :p", user.Name)
	}
	r.runRoulette(b, msg, bet)
	return nil
}

func (r *Roulette) runRoulette(b *bot.Bot, msg *common.Msg, points int) {
	user := &msg.User
	won := rand.Float32() >= 0.5
	if won {
		user.Points += points
		b.SayFormat(r.WinMessage, msg, points)
	} else {
		user.Points -= points
		b.SayFormat(r.LoseMessage, msg, points)
	}

}
