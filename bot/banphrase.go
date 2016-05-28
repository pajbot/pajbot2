package bot

import "strings"

func (bot *Bot) CheckForBanphrase(msg Msg) Action {
	a := &Action{}
	m := strings.ToLower(msg.Message)
	if strings.Contains(m, "www.com") {
		a.Response = bot.Ban(msg.Username, 10, "bad link")
	}
	if a.Response != "" {
		a.Match = true
	}
	return *a
}
