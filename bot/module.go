package bot

import (
	"log"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
)

/*
Module xD
*/
type Module interface {
	Check(bot *Bot, msg *common.Msg, action *Action) error
	// just pass in the bot so the module has access to everything, not just sql
	Init(bot *Bot) (id string, enabled bool)
	DeInit(bot *Bot)
	GetState() *basemodule.BaseModule
}

// GetModule returns the instance of the module
func (b *Bot) GetModule(moduleName string) Module {
	var state *basemodule.BaseModule

	for _, module := range b.AllModules {
		state = module.GetState()
		if state.ID == moduleName {
			return module
		}
	}

	return nil
}

// EnableModule enables the given module in the bot
func (b *Bot) EnableModule(module Module) {
	module.Init(b)
	module.GetState().SetEnabled()

	b.EnabledModules = append(b.EnabledModules, module)

	module.GetState().SaveState(b.Redis, b.Channel.Name)
}

// DisableModule disabled the given module in the bot
func (b *Bot) DisableModule(module Module) {
	module.DeInit(b)
	module.GetState().SetDisabled()

	moduleIndex := -1

	for i, sModule := range b.EnabledModules {
		if module == sModule {
			moduleIndex = i
			break
		}
	}

	if moduleIndex != -1 {
		b.EnabledModules = append(b.EnabledModules[:moduleIndex], b.EnabledModules[moduleIndex+1:]...)
	} else {
		log.Println("Something went wrong when disabling module")
		return
	}

	module.GetState().SaveState(b.Redis, b.Channel.Name)
}
