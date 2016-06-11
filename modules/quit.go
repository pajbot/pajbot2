package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

// Quit xD
type Quit struct {
}

// Ensure the module implements the interface properly
var _ Module = (*Quit)(nil)

// Check xD
func (module *Quit) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	m := strings.Split(msg.Message, " ")
	trigger := strings.ToLower(m[0])

	// TODO: Make sure the User level is the right level
	if (trigger == "!quit" || trigger == "!exit") && (msg.User.Name == "nuuls" || msg.User.Name == "pajlada") {
		b.Quit <- "ayy lmao something bad happened xD"
	}
	return nil
}
