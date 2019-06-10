package modules

import (
	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("action_performer", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:               "action_performer",
			name:             "Action performer",
			enabledByDefault: true,
			priority:         500000,

			maker: newActionPerformer,
		}
	})
}

type actionPerformer struct {
	base
}

func newActionPerformer(b base) pkg.Module {
	return &actionPerformer{
		base: b,
	}
}

func (m actionPerformer) OnMessage(event pkg.MessageEvent) pkg.Actions {
	// FIXME
	// return action.Do()
	return nil
}
