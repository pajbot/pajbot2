package modules

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/points"
	"github.com/pajlada/pajbot2/sqlmanager"
)

// Points module
type Points struct {
	Roulette *points.Roulette
}

var _ Module = (*Points)(nil)

func (module *Points) Init(sql *sqlmanager.SQLManager) {
	module.Roulette = &points.Roulette{
		WinMessage:  "$(source) won %d points in roulette and now has $(source.points) points VisLaud",
		LoseMessage: "$(source) lost %d points in roulette and now has $(source.points) LUL",
	}
}

// Check xD
func (module *Points) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if !strings.HasPrefix(msg.Text, "!") {
		return nil
	}
	m := strings.ToLower(msg.Text)
	spl := strings.Split(m, " ")
	trigger := spl[0]
	var args []string
	if len(spl) > 1 {
		args = spl[1:]
	}

	// using pts to not trigger other bots
	switch trigger {
	case "!givepts":
		err := points.GivePoints(b, &msg.User, args)
		if err != nil {
			b.Say(fmt.Sprint(err))
		}
	case "!pts":
		msg.Args = args
		b.SaySafe(b.Format("$(user.name) has $(user.points) points KKaper", msg))
	case "!roul":
		err := module.Roulette.Run(b, msg, args)
		if err != nil {
			b.Say(fmt.Sprint(err))
		}
	}

	return nil
}
