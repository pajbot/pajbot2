package pajbot

// Bot is one instance of the pajbot bot.
// One instance of the pajbot bot can run multiple bots and join multiple chats.
type Bot struct {
}

func New() Bot {
	return Bot{}
}
