package bot

import "fmt"

func NewBot(cfg BotConfig) {
	bot := &Bot{}
	bot.Read = cfg.Readchan
	bot.Send = cfg.Sendchan
	bot.Channel = cfg.Channel
	bot.init()
}

func (bot *Bot) init() {
	fmt.Printf("new bot in %s\n", bot.Channel)
	for {
		m := <-bot.Read
		fmt.Printf("#%s %s : %s\n", m.Channel, m.Username, m.Message)
		if m.MessageType == "sub" {
			fmt.Printf("%s subbed for %d months in a row\n", m.Username, m.Length)
		}
		go bot.Handle(m)
	}
}

func (bot *Bot) Say(message string) {
	m := fmt.Sprintf("PRIVMSG #%s : %s", bot.Channel, message)
	bot.Send <- m
}
