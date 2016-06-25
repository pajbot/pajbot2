package modules

import (
	"fmt"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/filter"
	"github.com/pajlada/pajbot2/sqlmanager"
)

/*
Banphrase xD
*/
type Banphrase struct {
	Filters []filter.Filter
}

// Ensure the module implements the interface properly
var _ Module = (*Banphrase)(nil)

// Init xD
func (module *Banphrase) Init(sql *sqlmanager.SQLManager) {
	module.Filters = []filter.Filter{
		&filter.Link{},
		&filter.Length{},
	}
	// xD
}

// Check xD
func (module *Banphrase) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	a := filter.BanAction{}
	m := msg.Message // TODO: confusables, github.com/FiloSottile/tr39-confusables this is kinda shitty
	for _, f := range module.Filters {
		f.Run(m, msg, &a)
		if a.Matched {
			action.Response = b.Ban(msg.User.Name, a.Level, a.Reason)
		}
		if msg.Length > 400 {
			action.Response = b.Ban(msg.User.Name, 5, fmt.Sprintf("msg too long %d", msg.Length))
		}
	}
	return nil
}
