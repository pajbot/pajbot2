package bot

import (
	"fmt"
	"strings"

	"github.com/pajlada/pajbot2/common"
)

/*
Say sends a PRIVMSG to the bots given channel
*/
func (b *Bot) Say(message string) {
	m := fmt.Sprintf("PRIVMSG #%s :%s ", b.Channel.Name, message)
	b.Send <- m
}

/*
Sayf sends a formatted PRIVMSG to the bots given channel
*/
func (b *Bot) Sayf(format string, a ...interface{}) {
	b.Say(fmt.Sprintf(format, a...))
}

/*
Mention sends a message with a pre-decided format:
@Username: message
*/
func (b *Bot) Mention(user common.User, message string) {
	b.Sayf("@%s: %s", user.Name, message)
}

/*
Mentionf sends a formatted message with a pre-decided format:
@Username: formatted message
*/
func (b *Bot) Mentionf(user common.User, format string, a ...interface{}) {
	b.Mention(user, fmt.Sprintf(format, a...))
}

/*
SayFormat sends a formatted and safe message to the bots channel
*/
func (b *Bot) SayFormat(line string, msg *common.Msg, a ...interface{}) {
	b.SaySafef(b.Format(line, msg), a...)
}

/*
SaySafef sends a formatted PRIVMSG to the bots given channel
*/
func (b *Bot) SaySafef(format string, a ...interface{}) {
	b.SaySafe(fmt.Sprintf(format, a...))
}

/*
SaySafe allows only harmless irc commands,
this should be used for commands added by users
*/
func (b *Bot) SaySafe(message string) {
	if !strings.HasPrefix(message, "/") && !strings.HasPrefix(message, ".") {
		b.Say(message)
		return
	}
	m := strings.Split(message, " ")
	cmd := m[0][1:] // remove "." or "/"
	switch cmd {
	case "me":
	case "timeout":
	case "unban":
	case "subscribers":
	case "subscribersoff":
	case "emoteonly":
	case "emoteonlyoff":
	default:
		message = " " + message
	}
	b.Say(message)
}
