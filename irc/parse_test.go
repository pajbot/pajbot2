package irc_test

import (
	"testing"

	"github.com/nuuls/pajbot2/bot"
	"github.com/nuuls/pajbot2/irc"
	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	var messageTests = []struct {
		input    string
		expected bot.Msg
	}{
		{
			input:    "@badges=subscriber/1;color=#FF468F;display-name=Ampzyh;emotes=;mod=0;room-id=11148817;subscriber=1;turbo=0;user-id=40910607;user-type= :ampzyh!ampzyh@ampzyh.tmi.twitch.tv PRIVMSG #pajlada :LIKE THAT SOn",
			expected: bot.Msg{Color: "#FF468F", Displayname: "Ampzyh", Emotes: []bot.Emote{}, Mod: false, Subscriber: true, Turbo: false, Usertype: "", Username: "ampzyh", Channel: "pajlada", Message: "LIKE THAT SOn", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    "@badges=;color=;display-name=alfredmechinisme;emotes=29:3-8;mod=0;room-id=11148817;subscriber=0;turbo=0;user-id=118459065;user-type= :alfredmechinisme!alfredmechinisme@alfredmechinisme.tmi.twitch.tv PRIVMSG #pajlada :Ha MVGame",
			expected: bot.Msg{Color: "", Displayname: "alfredmechinisme", Emotes: []bot.Emote{bot.Emote{EmoteType: "twitch", ID: "29", Name: "", Pos: []string(nil), Count: 1}}, Mod: false, Subscriber: false, Turbo: false, Usertype: "", Username: "alfredmechinisme", Channel: "pajlada", Message: "Ha MVGame", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    "@badges=subscriber/1;color=#FF468F;display-name=Ampzyh;emotes=;mod=0;room-id=11148817;subscriber=1;turbo=0;user-id=40910607;user-type= :ampzyh!ampzyh@ampzyh.tmi.twitch.tv PRIVMSG #pajlada :OH HOHO GANGIUNG UP",
			expected: bot.Msg{Color: "#FF468F", Displayname: "Ampzyh", Emotes: []bot.Emote{}, Mod: false, Subscriber: true, Turbo: false, Usertype: "", Username: "ampzyh", Channel: "pajlada", Message: "OH HOHO GANGIUNG UP", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    "@badges=broadcaster/1,subscriber/1;color=#CC44FF;display-name=pajlada;emotes=;mod=1;room-id=11148817;subscriber=1;turbo=0;user-id=11148817;user-type=mod :pajlada!pajlada@pajlada.tmi.twitch.tv PRIVMSG #pajlada :!ping",
			expected: bot.Msg{Color: "#CC44FF", Displayname: "pajlada", Emotes: []bot.Emote{}, Mod: true, Subscriber: true, Turbo: false, Usertype: "mod", Username: "pajlada", Channel: "pajlada", Message: "!ping", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    "@badges=moderator/1,subscriber/1;color=#2E8B57;display-name=pajbot;emotes=;mod=1;room-id=11148817;subscriber=1;turbo=0;user-id=82008718;user-type=mod :pajbot!pajbot@pajbot.tmi.twitch.tv PRIVMSG #pajlada :pajbot has been online for 6 hours and 42 minutes",
			expected: bot.Msg{Color: "#2E8B57", Displayname: "pajbot", Emotes: []bot.Emote{}, Mod: true, Subscriber: true, Turbo: false, Usertype: "mod", Username: "pajbot", Channel: "pajlada", Message: "pajbot has been online for 6 hours and 42 minutes", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    "@badges=;color=#00FF7F;display-name=Skyrisenow;emotes=;mod=0;room-id=11148817;subscriber=0;turbo=0;user-id=102262699;user-type= :skyrisenow!skyrisenow@skyrisenow.tmi.twitch.tv PRIVMSG #pajlada :!stats ruined",
			expected: bot.Msg{Color: "#00FF7F", Displayname: "Skyrisenow", Emotes: []bot.Emote{}, Mod: false, Subscriber: false, Turbo: false, Usertype: "", Username: "skyrisenow", Channel: "pajlada", Message: "!stats ruined", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    "@badges=turbo/1,subscriber/1;color=#0B8D6F;display-name=nuuls;emotes=3287:0-4;mod=0;room-id=11148817;subscriber=1;turbo=1;user-id=100229878;user-type= :nuuls!nuuls@nuuls.tmi.twitch.tv PRIVMSG #pajlada :MiniK",
			expected: bot.Msg{Color: "#0B8D6F", Displayname: "nuuls", Emotes: []bot.Emote{bot.Emote{EmoteType: "twitch", ID: "3287", Name: "", Pos: []string(nil), Count: 1}}, Mod: false, Subscriber: true, Turbo: true, Usertype: "", Username: "nuuls", Channel: "pajlada", Message: "MiniK", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    "@badges=broadcaster/1,subscriber/1;color=#CC44FF;display-name=pajlada;emotes=12:13-14;mod=1;room-id=11148817;subscriber=1;turbo=0;user-id=11148817;user-type=mod :pajlada!pajlada@pajlada.tmi.twitch.tv PRIVMSG #pajlada :ACTION MEME-MESSAGE :P",
			expected: bot.Msg{Color: "#CC44FF", Displayname: "pajlada", Emotes: []bot.Emote{bot.Emote{EmoteType: "twitch", ID: "12", Name: "", Pos: []string(nil), Count: 1}}, Mod: true, Subscriber: true, Turbo: false, Usertype: "mod", Username: "pajlada", Channel: "pajlada", Message: "ACTION MEME-MESSAGE :P", MessageType: "privmsg", Me: false, Length: 0},
		},
		{
			input:    ":old_riotgems!old_riotgems@old_riotgems.tmi.twitch.tv PRIVMSG #pajlada :\\ pajaHappy / elevator \\ pajaCool /",
			expected: bot.Msg{Color: "", Displayname: "", Emotes: []bot.Emote(nil), Mod: false, Subscriber: false, Turbo: false, Usertype: "", Username: "old_riotgems", Channel: "pajlada", Message: "\\ pajaHappy / elevator \\ pajaCool /", MessageType: "privmsg", Me: false, Length: 0},
		},
	}

	for _, tt := range messageTests {
		res := irc.Parse(tt.input)

		assert.Equal(t, tt.expected, res)
	}
}
