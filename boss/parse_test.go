package boss

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
		/*
			{
				input:    "@badges=subscriber/1;color=#FF468F;display-name=Ampzyh;emotes=;mod=0;room-id=11148817;subscriber=1;turbo=0;user-id=40910607;user-type= :ampzyh!ampzyh@ampzyh.tmi.twitch.tv PRIVMSG #pajlada :LIKE THAT SOn",
				expected: common.Msg{Color: "#FF468F", Displayname: "Ampzyh", Emotes: []common.Emote{}, Mod: false, Subscriber: true, Turbo: false, Usertype: "", Username: "ampzyh", Channel: "pajlada", Message: "LIKE THAT SOn", MessageType: "privmsg", Me: false, Length: 0},
			},
			{
				input:    "@badges=;color=;display-name=alfredmechinisme;emotes=29:3-8;mod=0;room-id=11148817;subscriber=0;turbo=0;user-id=118459065;user-type= :alfredmechinisme!alfredmechinisme@alfredmechinisme.tmi.twitch.tv PRIVMSG #pajlada :Ha MVGame",
				expected: common.Msg{Color: "", Displayname: "alfredmechinisme", Emotes: []common.Emote{common.Emote{EmoteType: "twitch", ID: "29", Name: "", Pos: []string(nil), Count: 1}}, Mod: false, Subscriber: false, Turbo: false, Usertype: "", Username: "alfredmechinisme", Channel: "pajlada", Message: "Ha MVGame", MessageType: "privmsg", Me: false, Length: 0},
			},
			{
				input:    "@badges=subscriber/1;color=#FF468F;display-name=Ampzyh;emotes=;mod=0;room-id=11148817;subscriber=1;turbo=0;user-id=40910607;user-type= :ampzyh!ampzyh@ampzyh.tmi.twitch.tv PRIVMSG #pajlada :OH HOHO GANGIUNG UP",
				expected: common.Msg{Color: "#FF468F", Displayname: "Ampzyh", Emotes: []common.Emote{}, Mod: false, Subscriber: true, Turbo: false, Usertype: "", Username: "ampzyh", Channel: "pajlada", Message: "OH HOHO GANGIUNG UP", MessageType: "privmsg", Me: false, Length: 0},
			},
			{
				input:    "@badges=moderator/1,subscriber/1;color=#2E8B57;display-name=pajbot;emotes=;mod=1;room-id=11148817;subscriber=1;turbo=0;user-id=82008718;user-type=mod :pajbot!pajbot@pajbot.tmi.twitch.tv PRIVMSG #pajlada :pajbot has been online for 6 hours and 42 minutes",
				expected: common.Msg{Color: "#2E8B57", Displayname: "pajbot", Emotes: []common.Emote{}, Mod: true, Subscriber: true, Turbo: false, Usertype: "mod", Username: "pajbot", Channel: "pajlada", Message: "pajbot has been online for 6 hours and 42 minutes", MessageType: "privmsg", Me: false, Length: 0},
			},
			{
				input:    "@badges=;color=#00FF7F;display-name=Skyrisenow;emotes=;mod=0;room-id=11148817;subscriber=0;turbo=0;user-id=102262699;user-type= :skyrisenow!skyrisenow@skyrisenow.tmi.twitch.tv PRIVMSG #pajlada :!stats ruined",
				expected: common.Msg{Color: "#00FF7F", Displayname: "Skyrisenow", Emotes: []common.Emote{}, Mod: false, Subscriber: false, Turbo: false, Usertype: "", Username: "skyrisenow", Channel: "pajlada", Message: "!stats ruined", MessageType: "privmsg", Me: false, Length: 0},
			},
			{
				input:    "@badges=turbo/1,subscriber/1;color=#0B8D6F;display-name=nuuls;emotes=3287:0-4;mod=0;room-id=11148817;subscriber=1;turbo=1;user-id=100229878;user-type= :nuuls!nuuls@nuuls.tmi.twitch.tv PRIVMSG #pajlada :MiniK",
				expected: common.Msg{Color: "#0B8D6F", Displayname: "nuuls", Emotes: []common.Emote{common.Emote{EmoteType: "twitch", ID: "3287", Name: "", Pos: []string(nil), Count: 1}}, Mod: false, Subscriber: true, Turbo: true, Usertype: "", Username: "nuuls", Channel: "pajlada", Message: "MiniK", MessageType: "privmsg", Me: false, Length: 0},
			},
		*/
		{
			input: "@badges=broadcaster/1,subscriber/1;color=#CC44FF;display-name=pajlada;emotes=12:13-14;mod=1;room-id=11148817;subscriber=1;turbo=0;user-id=11148817;user-type=mod :pajlada!pajlada@pajlada.tmi.twitch.tv PRIVMSG #pajlada :ACTION MEME-MESSAGE :P",
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
				Message: "ACTION MEME-MESSAGE :P",
				Channel: "pajlada",
				Type:    common.MsgPrivmsg,
				Length:  0,
				Me:      false,
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
				Message: "!ping",
				Channel: "pajlada",
				Type:    common.MsgPrivmsg,
				Length:  0,
				Me:      false,
				Emotes:  []common.Emote{},
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
				Message: "\\ pajaHappy / elevator \\ pajaCool /",
				Channel: "pajlada",
				Type:    common.MsgPrivmsg,
				Length:  0,
				Me:      false,
				Emotes:  []common.Emote(nil),
			},
		},
	}

	for _, tt := range messageTests {
		p := &parse{}
		res := p.Parse(tt.input)

		assert.Equal(t, tt.expected, res)
	}
}
