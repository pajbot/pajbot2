package bot

import "github.com/pajlada/pajbot2/common"

/*
Module xD
*/
type Module interface {
	Check(bot *Bot, msg *common.Msg, action *Action) error
}
