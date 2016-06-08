package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
)

/*
Banphrase xD
*/
type Banphrase struct {
}

// Ensure the module implements the interface properly
var _ Module = (*Banphrase)(nil)

// Check xD
func (module *Banphrase) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	m := strings.ToLower(msg.Message)
	if strings.Contains(m, "www.com") {
		action.Response = b.Ban(msg.User.Name, 10, "bad link")
		action.Stop = true
	}
	return nil
}
