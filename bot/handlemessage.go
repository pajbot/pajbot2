package bot

type Action struct {
	Match    bool
	Response string
}

/*
Handle attempts to handle the given message
*/
func (bot *Bot) Handle(msg Msg) {
	action := bot.CheckForBanphrase(msg)
	if action.Match {
		bot.Say(action.Response)
		return
	}
	action = bot.CheckForCommand(msg)
	if action.Match {
		bot.Say(action.Response)
		return
	}

}
