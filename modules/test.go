package modules

import (
	"io/ioutil"
	"net/http"
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
	if msg.User.Level > 1000 {
		m := strings.Split(msg.Text, " ")
		if m[0] == "!follow" {
			b.Twitter.Follow(m[1])
			b.Sayf("now streaming %s's timeline", m[1])
		}
	}
	if msg.User.Level > 1000 {
		m := strings.Split(msg.Text, " ")
		if m[0] == "!lasttweet" {
			tweet := b.Twitter.LastTweetString(m[1])
			b.Sayf("last tweet from %s ", tweet)
		}
	}

	if msg.User.Level > 1000 {
		m := strings.Split(msg.Text, " ")
		if m[0] == "!join" {
			b.Join <- m[1]
		}
	}
	if msg.User.Level > 1000 {
		m := strings.Split(msg.Text, " ")
		if m[0] == "!part" {
			b.Join <- "PART " + m[1]
		}
	}

	if msg.User.Level > 1000 {
		m := strings.Split(msg.Text, " ")
		if m[0] == "!spam" {
			loops, err := strconv.ParseUint(m[1], 10, 64)
			if err != nil {
				b.Sayf("%v", err)
			}
			text := strings.Join(m[2:], " ")
			var i uint64
			for i < loops {
				b.Say(text)
				i++
			}
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
