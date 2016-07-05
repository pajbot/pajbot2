package modules

import (
	"strconv"
	"strings"

	"github.com/pajlada/pajbot2/apirequest"
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/web"
)

/*
Test xD
*/
type Test struct {
}

// Ensure the module implements the interface properly
var _ Module = (*Test)(nil)

// Check xD
func (module *Test) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if msg.User.Level > 1000 {
		m := strings.Split(msg.Text, " ")
		if m[0] == "!say" {
			b.SayFormat(msg.Text[5:], msg)
		}
	}
	if msg.User.Level > 1000 {
		m := strings.Split(msg.Text, " ")
		if m[0] == "!testapi" {
			log.Debug(apirequest.GetStream(m[1]))
		}
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
				Data: map[string]string{
					"text": msg.Text,
					"user": msg.User.DisplayName,
				},
			},
		}
		web.Hub.Broadcast(wsMessage)
	}
	r9k, slow, sub := msg.Tags["r9k"], msg.Tags["slow"], msg.Tags["subs-only"]
	switch msg.Type {
	case common.MsgRoomState:
		log.Debug("GOT MSG ROOMSTATE MESSAGE: %s", msg.Tags)
		if r9k != "" && slow != "" {
			// Initial channel join
			b.Sayf("initial join. state: r9k:%s, slow:%s, sub:%s", r9k, slow, sub)
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
	return nil
}
