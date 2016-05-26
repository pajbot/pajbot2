package bot

import "strings"

func (bot *Bot) Handle(msg Msg) error {
	m := strings.Split(msg.Message, " ")
	trigger := strings.ToLower(m[0])
	if trigger == "!xd" {
		bot.Say("pajaSWA")
	}
	return nil
}
