package modules

import (
	"log"

	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("goodbye", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "goodbye",
			name:  "Goodbye",
			maker: newGoodbye,

			enabledByDefault: false,
		}
	})
}

type goodbye struct {
	base
}

func newGoodbye(b base) pkg.Module {
	m := &goodbye{
		base: b,
	}

	conn, err := m.bot.Events().Listen("on_quit", func() error {
		go m.bot.Say("cya lol")
		return nil
	}, 100)
	if err != nil {
		log.Println("Error listening to on_quit event:", err)
		// FIXME
		// return err
	}

	m.connections = append(m.connections, conn)

	return m
}
