package modules

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/sqlmanager"
)

/*
Top xD
*/
type Top struct {
	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*Top)(nil)

func (module *Top) topPoints(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	const limit = 5
	const category = "points"
	users := b.Redis.Top(b.Channel.Name, category, limit)
	result := []string{}
	for _, u := range users {
		result = append(result, fmt.Sprintf("%s: %d", u.NameNoPing(), u.Points))
	}
	b.Sayf("Top %d %s: %s", limit, category, strings.Join(result, ", "))
}

// Init xD
func (module *Top) Init(sql *sqlmanager.SQLManager) {
	topPointsCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"toppts",
			},
		},
		Function: module.topPoints,
	}
	module.commandHandler.AddCommand(topPointsCommand)
}

// Check xD
func (module *Top) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}
