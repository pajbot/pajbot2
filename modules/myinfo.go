package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// MyInfo xD
type MyInfo struct {
	common.BaseModule
}

// Ensure the module implements the interface properly
var _ Module = (*MyInfo)(nil)

// Check xD
func (module *MyInfo) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	m := strings.Split(msg.Text, " ")
	trigger := strings.ToLower(m[0])

	if trigger == "!myinfo" {
		b.Sayf("ID: %d, username: %s, type: %s, level: %d",
			msg.User.ID, msg.User.DisplayName, msg.User.Type, msg.User.Level)
	}
	return nil
}

// Init xD
func (module *MyInfo) Init(bot *bot.Bot) (string, bool) {
	return "myinfo", isModuleEnabled(bot, "myinfo", true)
}

// DeInit xD
func (module *MyInfo) DeInit(b *bot.Bot) {

}
