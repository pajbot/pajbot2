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

func modulesLoad(b *bot.Bot) {
	// TODO(pajlada): Select which modules should be loaded
	//                via a redis json list or something
	b.Modules = []bot.Module{
		&modules.Banphrase{},
		&modules.Command{},
		&modules.Pyramid{},
		&modules.Quit{},
		&modules.SubAnnounce{},
		&modules.MyInfo{},
		&modules.Test{},
		&modules.Admin{},
		&modules.Points{},
		&modules.Top{},
		&modules.Raffle{},
		&modules.Bingo{},
	}

	// Initialize all loaded modules
	for _, module := range b.Modules {
		module.Init(b)
	}
}

// modulesReload unloads all loaded modules, then reloads all modules
// that should be enabled
func modulesReload(b *bot.Bot) {
	modulesUnload(b)
	modulesLoad(b)
}
