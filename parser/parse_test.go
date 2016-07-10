package parser

import (
	"testing"

	"github.com/pajlada/pajbot2/common"
	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	var messageTests = []struct {
		input    string
		expected common.Msg
	}{
		{
			input: "@badges=broadcaster/1,subscriber/1;color=#CC44FF;display-name=pajlada;emotes=12:13-14;mod=1;room-id=11148817;subscriber=1;turbo=0;user-id=11148817;user-type=mod :pajlada!pajlada@pajlada.tmi.twitch.tv PRIVMSG #pajlada :\u0001ACTION MEME-MESSAGE :P\u0001",
			expected: common.Msg{
				User: common.User{
					ID:           0,
					Name:         "pajlada",
					DisplayName:  "pajlada",
					Mod:          true,
					Sub:          true,
					Turbo:        false,
					ChannelOwner: true,
					Type:         "mod",
					Level:        0,
					Points:       0,
				},
				Text:    "MEME-MESSAGE :P",
				Channel: "pajlada",
				Type:    common.MsgPrivmsg,
				Me:      true,
				Emotes: []common.Emote{
					{
						Name:  ":P",
						ID:    "12",
						Type:  "twitch",
						SizeX: 28,
						SizeY: 28,
						IsGif: false,
						Count: 1,
					},
				},
				Tags: map[string]string{
					"badges":  "broadcaster/1,subscriber/1",
					"color":   "#CC44FF",
					"room-id": "11148817",
					"user-id": "11148817",
				},
			},
		},
		{
			input: "@badges=subscriber/1;color=#FF0000;display-name=JoeMoneyTV;emotes=;login=joemoneytv;mod=0;msg-id=resub;msg-param-months=6;room-id=11148817;subscriber=1;system-msg=JoeMoneyTV\\ssubscribed\\sfor\\s6\\smonths\\sin\\sa\\srow!;turbo=0;user-id=56871381;user-type= :tmi.twitch.tv USERNOTICE #pajlada :huehueheuheue",
			expected: common.Msg{
				User: common.User{
					ID:           0,
					Name:         "joemoneytv",
					DisplayName:  "JoeMoneyTV",
					Mod:          false,
					Sub:          true,
					Turbo:        false,
					ChannelOwner: false,
					Type:         "",
					Level:        0,
					Points:       0,
				},
				Text:    "huehueheuheue",
				Channel: "pajlada",
				Type:    common.MsgReSub,
				Me:      false,
				Emotes:  []common.Emote{},
				Tags: map[string]string{
					"badges":           "subscriber/1",
					"color":            "#FF0000",
					"room-id":          "11148817",
					"user-id":          "56871381",
					"msg-param-months": "6",
					"msg-id":           "resub",
					"login":            "joemoneytv",
					"system-msg":       "JoeMoneyTV subscribed for 6 months in a row!",
				},
			},
		},
		{
			input: "@badges=subscriber/1;color=#FF0000;display-name=JoeMoneyTV;emotes=;login=joemoneytv;mod=0;msg-id=resub;msg-param-months=6;room-id=11148817;subscriber=1;system-msg=JoeMoneyTV\\ssubscribed\\sfor\\s6\\smonths\\sin\\sa\\srow!;turbo=0;user-id=56871381;user-type= :tmi.twitch.tv USERNOTICE #pajlada",
			expected: common.Msg{
				User: common.User{
					ID:           0,
					Name:         "joemoneytv",
					DisplayName:  "JoeMoneyTV",
					Mod:          false,
					Sub:          true,
					Turbo:        false,
					ChannelOwner: false,
					Type:         "",
					Level:        0,
					Points:       0,
				},
				Text:    "",
				Channel: "pajlada",
				Type:    common.MsgReSub,
				Me:      false,
				Emotes:  []common.Emote{},
				Tags: map[string]string{
					"badges":           "subscriber/1",
					"color":            "#FF0000",
					"room-id":          "11148817",
					"user-id":          "56871381",
					"msg-param-months": "6",
					"msg-id":           "resub",
					"login":            "joemoneytv",
					"system-msg":       "JoeMoneyTV subscribed for 6 months in a row!",
				},
			},
		},
		{
			input: "@badges=broadcaster/1,subscriber/1;color=#CC44FF;display-name=pajlada;emotes=;mod=1;room-id=11148817;subscriber=1;turbo=0;user-id=11148817;user-type=mod :pajlada!pajlada@pajlada.tmi.twitch.tv PRIVMSG #pajlada :!ping",
			expected: common.Msg{
				User: common.User{
					ID:           0,
					Name:         "pajlada",
					DisplayName:  "pajlada",
					Mod:          true,
					Sub:          true,
					Turbo:        false,
					ChannelOwner: true,
					Type:         "mod",
					Level:        0,
					Points:       0,
				},
				Text:    "!ping",
				Channel: "pajlada",
				Type:    common.MsgPrivmsg,
				Me:      false,
				Emotes:  []common.Emote{},
				Tags: map[string]string{
					"badges":  "broadcaster/1,subscriber/1",
					"color":   "#CC44FF",
					"room-id": "11148817",
					"user-id": "11148817",
				},
			},
		},
		{
			input: ":old_riotgems!old_riotgems@old_riotgems.tmi.twitch.tv PRIVMSG #pajlada :\\ pajaHappy / elevator \\ pajaCool /",
			expected: common.Msg{
				User: common.User{
					ID:          0,
					Name:        "old_riotgems",
					DisplayName: "",
					Mod:         false,
					Sub:         false,
					Turbo:       false,
					Type:        "",
					Level:       0,
					Points:      0,
				},
				Text:    "\\ pajaHappy / elevator \\ pajaCool /",
				Channel: "pajlada",
				Type:    common.MsgPrivmsg,
				Me:      false,
				Emotes:  []common.Emote(nil),
				Tags:    nil,
			},
		},
	}

	for _, tt := range messageTests {
		res := Parse(tt.input)

		assert.Equal(t, tt.expected, res)
	}
}
