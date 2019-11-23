package tusecommands

import (
	"errors"
	"fmt"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	mcommands "github.com/pajbot/pajbot2/pkg/modules/commands"
)

const id = "tusecommands"
const name = "Test module - Use commands module"

func init() {
	modules.Register(id, func() pkg.ModuleSpec {
		return modules.NewSpec(id, name, false, newModule)
	})
}

type module struct {
	mbase.Base
}

func (m *module) registerCommand() error {
	iCommandsModule, err := m.BotChannel().GetModule("commands")
	if err != nil {
		return err
	}

	commandsModule, ok := iCommandsModule.(*mcommands.CommandsModule)
	if !ok {
		return errors.New("this module is not a commands module wtf")
	}

	fmt.Println("got commands module:", commandsModule)

	return nil
}

func newModule(b mbase.Base) pkg.Module {
	m := &module{
		Base: b,
	}

	return m
}
