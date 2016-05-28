package bot

/*
Handle attempts to handle the given message
*/
func (bot *Bot) Handle(msg Msg) {
	action := &Action{}
	for _, module := range bot.Modules {
		module.Check(bot, &msg, action)

		if action.Response != "" {
			bot.Say(action.Response)
		}

		if action.Stop {
			return
		}
	}
}
