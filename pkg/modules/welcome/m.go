package welcome

import (
	"log"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	modules.Register("welcome", func() pkg.ModuleSpec {
		return modules.NewSpec("welcome", "Welcome", false, newWelcome)
	})
}

type welcome struct {
	mbase.Base
}

func newWelcome(b *mbase.Base) pkg.Module {
	m := &welcome{
		Base: *b,
	}

	// FIXME
	m.Initialize()

	return m
}

func (m *welcome) Initialize() {
	err := m.Listen("on_join", func() error {
		go m.BotChannel().Say("pb2 joined")
		return nil
	}, 100)
	if err != nil {
		// FIXME
		log.Println("ERROR LISTENING TO ON JOIN XD")
	}
}
