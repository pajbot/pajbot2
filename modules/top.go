package modules

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
)

/*
Top xD
*/
type Top struct {
	common.BaseModule
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

func (module *Top) topSpammerOnline(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	const limit = 5
	const category = "online_message_count"
	const categoryTitle = "online chatter"
	users := b.Redis.Top(b.Channel.Name, category, limit)
	result := []string{}
	for _, u := range users {
		result = append(result, fmt.Sprintf("%s: %d", u.NameNoPing(), u.OnlineMessageCount))
	}
	b.Sayf("Top %d %s: %s", limit, categoryTitle, strings.Join(result, ", "))
}

func (module *Top) topSpammerOffline(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	const limit = 5
	const category = "offline_message_count"
	const categoryTitle = "offline chatter"
	users := b.Redis.Top(b.Channel.Name, category, limit)
	result := []string{}
	for _, u := range users {
		result = append(result, fmt.Sprintf("%s: %d", u.NameNoPing(), u.OfflineMessageCount))
	}
	b.Sayf("Top %d %s: %s", limit, categoryTitle, strings.Join(result, ", "))
}
func (module *Top) topSpammerTotal(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	const limit = 5
	const category = "total_message_count"
	const categoryTitle = "chatter"
	users := b.Redis.Top(b.Channel.Name, category, limit)
	result := []string{}
	for _, u := range users {
		result = append(result, fmt.Sprintf("%s: %d", u.NameNoPing(), u.TotalMessageCount))
	}
	b.Sayf("Top %d %s: %s", limit, categoryTitle, strings.Join(result, ", "))
}

// Init xD
func (module *Top) Init(bot *bot.Bot) (string, bool) {
	topPointsCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"toppts",
			},
		},
		Function: module.topPoints,
	}
	module.commandHandler.AddCommand(topPointsCommand)
	topSpammerOnlineCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"online",
			},
		},
		Function: module.topSpammerOnline,
	}
	topSpammerOfflineCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"offline",
			},
		},
		Function: module.topSpammerOffline,
	}
	topSpammerTotalCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"all",
			},
		},
		Function: module.topSpammerTotal,
	}
	topSpammerCommand := &command.NestedCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"topchatter",
				"topchatters",
				"topspammer",
				"topspammers",
			},
		},
		Commands: []command.Command{
			topSpammerOnlineCommand,
			topSpammerOfflineCommand,
			topSpammerTotalCommand,
		},
		DefaultCommand:  topSpammerTotalCommand,
		FallbackCommand: topSpammerTotalCommand,
	}
	module.commandHandler.AddCommand(topSpammerCommand)

	return "top", isModuleEnabled("top")
}

// DeInit xD
func (module *Top) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Top) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}
