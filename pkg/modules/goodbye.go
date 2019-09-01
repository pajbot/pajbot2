package modules

import (
	"log"

	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	Register("goodbye", func() pkg.ModuleSpec {
		return &Spec{
			id:    "goodbye",
			name:  "Goodbye",
			maker: newGoodbye,

			enabledByDefault: false,
		}
	})
}

type goodbye struct {
	mbase.Base
}

func newGoodbye(b mbase.Base) pkg.Module {
	m := &goodbye{
		Base: b,
	}

	err := m.Listen("on_quit", func() error {
		go m.BotChannel().Say("cya lol")
		return nil
	}, 100)
	if err != nil {
		log.Println("Error listening to on_quit event:", err)
		// FIXME
		// return err
	}

	return m
}
