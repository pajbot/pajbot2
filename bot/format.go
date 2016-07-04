package bot

import (
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/format"
)

// Format formats the given line xD
func (bot *Bot) Format(line string, msg *common.Msg) string {
	// catch all errors until we have proper error handling
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()
	fmtline, rawCommands := format.ParseLine(line)
	for i := range rawCommands {
		format.ExecCommand(bot.Redis, &rawCommands[i], msg)
	}
	return format.RunCommands(fmtline, rawCommands)
}
