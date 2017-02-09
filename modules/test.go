package modules

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dankeroni/gotwitch"
	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/basemodule"
	"github.com/pajlada/pajbot2/helper"
	"github.com/pajlada/pajbot2/web"
)

/*
Test xD
*/
type Test struct {
	basemodule.BaseModule
	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*Test)(nil)

func cmdTest(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	m := helper.GetTriggersN(msg.Text, 1)

	if len(m) == 0 {
		// missing argument to !test
		return
	}

	switch m[0] {
	case "say":
		b.Say(strings.Join(m[1:], " "))
	case "lasttweet":
		if len(m) > 1 {
			tweet := b.Twitter.LastTweetString(m[1])
			b.Sayf("last tweet from %s ", tweet)
		} else {
			b.Say("Usage: !test lasttweet pajlada")
		}
	case "follow":
		if len(m) > 1 {
			b.Twitter.Follow(m[1])
			b.Sayf("now streaming %s's timeline", m[1])
		} else {
			b.Say("Usage: !test follow pajlada")
		}
	case "unfollow":
		b.Say("not implemented yet")
	case "api":
		if len(m) > 1 {
			apirequest.Twitch.GetStream(m[1],
				func(stream gotwitch.Stream) {
					b.Sayf("Stream info: %+v", stream)
				},
				func(statusCode int, statusMessage, errorMessage string) {
					b.Sayf("ERROR: %d", statusCode)
					b.Say(statusMessage)
					b.Say(errorMessage)
				}, func(err error) {
					b.Say("Internal error")
				})
		} else {
			b.Say("Usage: !test api pajlada")
		}
	case "resub":
		testMessage := `@badges=staff/1,broadcaster/1,turbo/1;color=#008000;display-name=TWITCH_UserName;emotes=;mod=0;msg-id=resub;msg-param-months=6;room-id=1337;subscriber=1;system-msg=TWITCH_UserName\shas\ssubscribed\sfor\s6\smonths!;login=twitch_username;turbo=1;user-id=1337;user-type=staff :tmi.twitch.tv USERNOTICE #%s :Great stream -- keep it up!`
		b.RawRead <- fmt.Sprintf(testMessage, b.Channel.Name)

	case "whisper":
		log.Debugf("WHISPER %s", msg.User.Name)
		b.Whisper(msg.User.Name, "TEST WHISPER")

	default:
		b.Sayf("Unhandled action %s", m[0])
		return
	}
}

// Init xD
func (module *Test) Init(bot *bot.Bot) (string, bool) {
	module.SetDefaults("test")
	module.EnabledDefault = true
	module.ParseState(bot.Redis, bot.Channel.Name)

	testCommand := command.FuncCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"test",
			},
			Level: 1000,
		},
		Function: cmdTest,
	}
	module.commandHandler.AddCommand(&testCommand)

	return "test", true
}

// DeInit xD
func (module *Test) DeInit(b *bot.Bot) {

}

// Check xD
func (module *Test) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if strings.HasPrefix(msg.Text, "!") {
		trigger := strings.Split(msg.Text, " ")[0]
		if strings.ToLower(trigger) == "!relaybroker" {
			req, err := http.Get("http://localhost:9002/stats")
			if err != nil {
				log.Error(err)
				return nil
			}
			bs, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Error(err)
			}
			b.SaySafe(string(bs))
		}
	}

	if len(msg.Emotes) > 0 {
		wsMessage := &web.WSMessage{
			MessageType: web.MessageTypeCLR,
			Payload: &web.Payload{
				Event: "emotes",
				Data: map[string]interface{}{
					"user":   msg.User.DisplayName,
					"emotes": msg.Emotes,
				},
			},
		}
		web.Hub.Broadcast(wsMessage)
	}

	if msg.Text == "abc" {
		wsMessage := &web.WSMessage{
			MessageType: web.MessageTypeDashboard,
			Payload: &web.Payload{
				Event: "xD",
			},
		}
		web.Hub.Broadcast(wsMessage)
	} else {
		wsMessage := &web.WSMessage{
			MessageType: web.MessageTypeDashboard,
			Payload: &web.Payload{
				Event: "chat",
				Data: map[string]interface{}{
					"text": msg.Text,
					"user": msg.User.DisplayName,
				},
			},
		}
		web.Hub.Broadcast(wsMessage)
	}
	switch msg.Type {
	case common.MsgTimeoutSuccess:
		// b.Sayf("MsgTimeoutSuccess triggered: %#v", msg.Tags)
	case common.MsgRoomState:
		log.Debug("GOT MSG ROOMSTATE MESSAGE: %s", msg.Tags)
		r9k, slow, sub := msg.Tags["r9k"], msg.Tags["slow"], msg.Tags["subs-only"]
		if r9k != "" && slow != "" {
			// Initial channel join
			//b.Sayf("initial join. state: r9k:%s, slow:%s, sub:%s", r9k, slow, sub)
			b.Say("MrDestructoid")
		} else {
			if r9k != "" {
				if r9k == "1" {
					b.Say("r9k on")
				} else {
					b.Say("r9k off")
				}
			} else if slow != "" {
				slowDuration, err := strconv.Atoi(slow)
				if err == nil {
					if slowDuration == 0 {
						b.Say("Slowmode off")
					} else {
						b.Sayf("Slowmode changed to %d seconds", slowDuration)
					}
				}
			} else if sub != "" {
				if sub == "1" {
					b.Say("submode on")
				} else {
					b.Say("submode off")
				}
			}
		}
	}

	return module.commandHandler.Check(b, msg, action)
}
