package modules

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
	"github.com/pajlada/pajbot2/web"
)

/*
BroadcastEmotes xD
*/
type BroadcastEmotes struct {
	basemodule.BaseModule
	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*BroadcastEmotes)(nil)

// Init xD
func (module *BroadcastEmotes) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("test")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	return "test", true
}

// DeInit xD
func (module *BroadcastEmotes) DeInit(b *bot.Bot) {

}

// Check xD
func (module *BroadcastEmotes) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if len(msg.Emotes) > 0 {
		wsMessage := &web.WSMessage{
			MessageType: web.MessageTypeCLR,
			Payload: &web.Payload{
				Event: "emotes",
				Data: map[string]interface{}{
					"user":   msg.User.DisplayName,
					"emotes": msg.Emotes,
				},
			},
		}
		web.Hub.Broadcast(wsMessage)
	}

	return nil
}
