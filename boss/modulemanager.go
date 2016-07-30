package boss

import (
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/modules"
)

func modulesUnload(b *bot.Bot) {
	// De-init all already-loaded modules
	for _, module := range b.Modules {
		module.DeInit(b)
	}

	b.Modules = nil
}

func modulesInit(b *bot.Bot) {
	// TODO(pajlada): Select which modules should be loaded
	//                via a redis json list or something
	b.AllModules = []bot.Module{
		&modules.Admin{},
		&modules.Banphrase{},
		&modules.Bingo{},
		&modules.Command{},
		&modules.MyInfo{},
		&modules.Points{},
		&modules.Pyramid{},
		&modules.Quit{},
		&modules.Raffle{},
		&modules.SubAnnounce{},
		&modules.Test{},
		&modules.Top{},
	}
}

func modulesLoad(b *bot.Bot) {
	// Initialize all loaded modules
	for _, module := range b.AllModules {
		module.Init(b)
	}

	b.Modules = nil

	for _, module := range b.AllModules {
		state := module.GetState()
		if state.IsEnabled() {
			log.Debugf("Enabling module %s", state.ID)
			b.Modules = append(b.Modules, module)
		} else {
			log.Debugf("Module %s will not be enabled", state.ID)
		}
	}
}

// modulesReload unloads all loaded modules, then reloads all modules
// that should be enabled
func modulesReload(b *bot.Bot) {
	modulesUnload(b)
	modulesLoad(b)
}

// Compile BaseModules and OptionalModules into the Modules slice
// This will be based on whether Init returned true or not
// Or maybe there should be another method. Like Enabled?
// Can we make it so all modules have the same thing?
// Or do we need to reimplement the wheel?
func modulesCompile(b *bot.Bot) {

}
