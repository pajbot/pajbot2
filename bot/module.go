package bot

import "github.com/pajlada/pajbot2/common"

/*
Module xD
*/
type Module interface {
	Check(bot *Bot, msg *common.Msg, action *Action) error
	// just pass in the bot so the module has access to everything, not just sql
	Init(bot *Bot) (id string, enabled bool)
	DeInit(bot *Bot)
	GetState() *common.BaseModule
	SetState(id string, enabled bool)
}
