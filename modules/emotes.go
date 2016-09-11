package modules

import (
	"sort"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
)

/*
Emotes xD
*/
type Emotes struct {
	basemodule.BaseModule
	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*Emotes)(nil)

func (module *Emotes) ffzEmotes(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	var activeEmotes []string

	for _, emote := range b.Channel.Emotes.FrankerFaceZ {
		activeEmotes = append(activeEmotes, emote.Name)
	}
	sort.Strings(activeEmotes)

	b.SaySafef("Active FFZ emotes: %s", strings.Join(activeEmotes, " "))

	if msg.User.Level >= 1000 {
		b.SaySafef("Last updated: %s (use !emotes reload ffz)", b.Channel.Emotes.FrankerFaceZLastUpdate)
	}
}

func (module *Emotes) bttvEmotes(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	var activeEmotes []string

	for _, emote := range b.Channel.Emotes.Bttv {
		activeEmotes = append(activeEmotes, emote.Name)
	}
	sort.Strings(activeEmotes)

	b.SaySafef("Active BTTV emotes: %s", strings.Join(activeEmotes, " "))

	if msg.User.Level >= 1000 {
		b.SaySafef("Last updated: %s (use !emotes reload bttv)", b.Channel.Emotes.BttvLastUpdate)
	}
}

func (module *Emotes) bttvGlobalEmotes(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	var activeEmotes []string

	for _, emote := range bot.GlobalEmotes.Bttv {
		activeEmotes = append(activeEmotes, emote.Name)
	}
	sort.Strings(activeEmotes)

	b.SaySafef("BTTV global emotes: %s", strings.Join(activeEmotes, " "))
}

// Init xD
func (module *Emotes) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("emotes")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	emotesBttvCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"bttv",
				"bettertwitchtv",
			},
		},
		Function: module.bttvEmotes,
	}
	emotesGlobalBttvCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"globalbttv",
				"globalbettertwitchtv",
			},
			Level: 1000,
		},
		Function: module.bttvGlobalEmotes,
	}

	emotesFfzCommand := &command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"ffz",
				"frankerfacez",
			},
		},
		Function: module.ffzEmotes,
	}

	emotesCommand := &command.NestedCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"emotes",
			},
		},
		Commands: []command.Command{
			emotesBttvCommand,
			emotesFfzCommand,
			emotesGlobalBttvCommand,
		},
	}

	module.commandHandler.AddCommand(emotesCommand)

	return "emotes", true
}

// DeInit xD
func (module *Emotes) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Emotes) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}
