package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// Quit xD
type Quit struct {
	common.BaseModule
}

// Ensure the module implements the interface properly
var _ Module = (*Quit)(nil)

// Init xD
func (module *Quit) Init(bot *bot.Bot) (string, bool) {
	// XXX: MOVE THIS
	return "ASD", true
}

// DeInit xD
func (module *Quit) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Quit) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	m := strings.Split(msg.Text, " ")
	trigger := strings.ToLower(m[0])

	// TODO: Make sure the User level is the right level
	if (trigger == "!quit" || trigger == "!exit") && (msg.User.Name == "nuuls" || msg.User.Name == "pajlada") {
		b.Quit <- "quit from command by " + msg.User.Name
	}
	return nil
}
