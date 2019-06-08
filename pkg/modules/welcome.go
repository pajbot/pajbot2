package modules

import (
	"log"

	"github.com/pajbot/pajbot2/pkg"
)

func init() {
	Register("welcome", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "welcome",
			name:  "Welcome",
			maker: newWelcome,

			enabledByDefault: false,
		}
	})
}

type welcome struct {
	base
}

func newWelcome(b base) pkg.Module {
	m := &welcome{
		base: b,
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *welcome) Initialize() {
	conn, err := m.bot.Events().Listen("on_join", func() error {
		go m.bot.Say("pb2 joined")
		return nil
	}, 100)
	if err != nil {
		// FIXME
		log.Println("ERROR LISTENING TO ON JOIN XD")
	}

	m.connections = append(m.connections, conn)
}
