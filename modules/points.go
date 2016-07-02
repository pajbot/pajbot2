package modules

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/points"
)

// Points module
type Points struct {
}

var _ Module = (*Points)(nil)

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
		b.Sayf("%s has %d points KKaper", msg.User.Name, msg.User.Points)
	case "!roul":
		err := points.Roulette(b, &msg.User, args)
		if err != nil {
			b.Say(fmt.Sprint(err))
		}
	}

	return nil
}
