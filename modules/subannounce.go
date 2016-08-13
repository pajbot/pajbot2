package modules

import (
	"fmt"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
)

/*
SubAnnounce xD
*/
type SubAnnounce struct {
	basemodule.BaseModule
}

// Ensure the module implements the interface properly
var _ Module = (*SubAnnounce)(nil)

// Init xD
func (module *SubAnnounce) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("sub-announce")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	return "sub-announce", true
}

// DeInit xD
func (module *SubAnnounce) DeInit(b *bot.Bot) {

}

// Check xD
func (module *SubAnnounce) Check(b *bot.Bot, m *common.Msg, action *bot.Action) error {
	if m.Type == common.MsgSub {
		action.Response = fmt.Sprintf("%s just subscribed! PogChamp", m.User.Name)
		action.Stop = true
	} else if m.Type == common.MsgReSub {
		action.Response = fmt.Sprintf("%s just subscribed for %d months in a row! PogChamp", m.User.Name, 1337)
		action.Stop = true
	}
	return nil
}
