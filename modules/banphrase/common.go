package banphrase

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/modules"
)

/*
Banphrase xD
*/
type Banphrase struct {
}

// Ensure the module implements the interface properly
var _ module.Module = (*Banphrase)(nil)

// Check xD
func (module *Banphrase) Check(b *bot.Bot, msg *bot.Msg, action *bot.Action) error {
	m := strings.ToLower(msg.Message)
	if strings.Contains(m, "www.com") {
		action.Response = b.Ban(msg.Username, 10, "bad link")
		action.Stop = true
	}
	return nil
}
