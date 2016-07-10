package parser

import (
	"strconv"
	"strings"

	"github.com/pajlada/pajbot2/common"
)

/*
Parse parses an IRC message into a more readable bot.Msg
*/
func Parse(line string) common.Msg {
	m := &common.Msg{
		User: common.User{},
		Type: common.MsgUnknown,
	}

	// msg is the string we will keep working on/reducing as we parse things
	msg := line

	var splitLine []string

	// The message starts with @, that means there are IRCv3 tags available to parse
	if strings.HasPrefix(line, "@") {
		splitLine = strings.SplitN(msg, " ", 2)
		parseTags(m, splitLine[0][1:])
		msg = splitLine[1]
	}

	// Parse source
	splitLine = strings.SplitN(msg, " ", 2)
	parseSource(m, splitLine[0])
	msg = splitLine[1]

	// Parse message type
	splitLine = strings.SplitN(msg, " ", 2)
	parseMsgType(m, splitLine[0])
	msg = splitLine[1]

	if m.Type == common.MsgUnknown {
		m.Type = common.MsgThrowAway
		return *m
	}

	splitLine = strings.SplitN(msg, " ", 2)
	parseChannel(m, splitLine[0])

	if len(splitLine) == 2 {
		msg = splitLine[1]

		// Parse message text + msg type (if it's a /me message or not)
		parseText(m, msg)

		if m.User.Name == "twitchnotify" {
			if !strings.Contains(m.Text, " to ") && !strings.Contains(m.Text, " while ") {
				parseNewSub(m)
			}
		}
	}

	// If the destination of the message is the same as the username,
	// then we tag the user as the channel owner. This will automatically
	// give him access to broadcaster commands
	if m.Channel == m.User.Name {
		m.User.ChannelOwner = true
	}

	if m.Tags != nil {
		// Parse tags further, such as the msg-id value for determinig the msg type
		parseExtendedTags(m)
	}

	return *m
}

func parseTwitchEmotes(m *common.Msg, emotetag string) {
	// TODO: Parse more emote information (bttv (and ffz?), name, size, isGif)
	// will we done by a module in the bot itself
	m.Emotes = make([]common.Emote, 0)
	if emotetag == "" {
		return
	}
	emoteSlice := strings.Split(emotetag, "/")
	for i := range emoteSlice {
		spl := strings.Split(emoteSlice[i], ":")
		id := spl[0]
		e := &common.Emote{}
		e.Type = "twitch"
		e.Name = getEmoteName(m, spl[1])
		e.ID = id
		// 28 px should be fine for twitch emotes
		e.SizeX = 28
		e.SizeY = 28
		e.Count = strings.Count(emoteSlice[i], "-")
		m.Emotes = append(m.Emotes, *e)
	}
}

func getEmoteName(m *common.Msg, pos string) string {
	pos = strings.Split(pos, ",")[0]
	spl := strings.Split(pos, "-")
	start, _ := strconv.Atoi(spl[0])
	end, _ := strconv.Atoi(spl[1])
	runes := []rune(m.Text)
	name := runes[start : end+1]
	return string(name)
}

func parseTagValues(m *common.Msg) {
	// TODO: Parse id and color
	// color and id is pretty useless imo
	if m.Tags["display-name"] == "" {
		m.User.DisplayName = m.User.Name
	} else {
		m.User.DisplayName = m.Tags["display-name"]
	}
	delete(m.Tags, "display-name")
	m.User.Type = m.Tags["user-type"]
	delete(m.Tags, "user-type")
	// fucking linter
	one := "1"
	if m.Tags["turbo"] == one {
		m.User.Turbo = true
	}
	delete(m.Tags, "turbo")
	if m.Tags["mod"] == one {
		m.User.Mod = true
	}
	delete(m.Tags, "mod")
	if m.Tags["subscriber"] == one {
		m.User.Sub = true
	}
	delete(m.Tags, "subscriber")
}

func parseExtendedTags(m *common.Msg) {
	// Parse twitch emotes from the "emotes" tag
	parseTwitchEmotes(m, m.Tags["emotes"])
	delete(m.Tags, "emotes")

	switch m.Tags["msg-id"] {
	case "resub":
		m.Type = common.MsgReSub

	case "subs_on":
		m.Type = common.MsgSubsOn

	case "subs_off":
		m.Type = common.MsgSubsOff

	case "slow_on":
		// Slow mode duration is found in the tag slow_duration
		m.Type = common.MsgSlowOn

	case "slow_off":
		m.Type = common.MsgSlowOff

	case "r9k_on":
		m.Type = common.MsgR9kOn

	case "r9k_off":
		m.Type = common.MsgR9kOff

	case "host_on":
		// Host target can be found in target_channel tag
		m.Type = common.MsgHostOn

	case "host_off":
		m.Type = common.MsgHostOff

	case "":
		break

	default:
		m.Type = common.MsgUnknown
	}

	if m.Tags["login"] != "" {
		m.User.Name = m.Tags["login"]
	}
}

/*
XXX: Should user properties stay at their zero value when there are no tags? Do we even care about this scenario?
*/
func parseTags(m *common.Msg, msg string) {
	m.Tags = make(map[string]string)
	// IRCv3-tags are separated by semicolons
	for _, tagValue := range strings.Split(msg, ";") {
		spl := strings.Split(tagValue, "=")
		k := spl[0]
		v := strings.Replace(spl[1], "\\s", " ", -1)
		m.Tags[k] = v
	}

	parseTagValues(m)

}

func parseSource(m *common.Msg, msg string) {
	if strings.HasPrefix(msg, ":") {
		msg = msg[1:]
	}
	// Check if the source is a user
	userSepPos := strings.Index(msg, "!")
	hostSepPos := strings.Index(msg, "@")
	if userSepPos > -1 && hostSepPos > -1 && userSepPos < hostSepPos {
		// A valid user address is found!
		m.User.Name = msg[0:userSepPos]
	}
}

func parseMsgType(m *common.Msg, msg string) {
	switch msg {
	case "PRIVMSG":
		m.Type = common.MsgPrivmsg

	case "WHISPER":
		m.Type = common.MsgWhisper

	case "USERNOTICE":
		m.Type = common.MsgUserNotice

	case "NOTICE":
		m.Type = common.MsgNotice

	case "ROOMSTATE":
		m.Type = common.MsgRoomState
	}
}

func parseChannel(m *common.Msg, msg string) {
	m.Channel = strings.Replace(msg[1:], "#", "", 0)
}

func parseText(m *common.Msg, msg string) {
	m.Text = msg[1:]

	// figure out whether the message is an ACTION or not
	getAction(m)
}

// regex in 2016 LUL
func getAction(m *common.Msg) {
	if strings.HasPrefix(m.Text, "\u0001ACTION ") && strings.HasSuffix(m.Text, "\u0001") {
		m.Me = true
		msg := m.Text
		msg = strings.Replace(msg, "\u0001ACTION ", "", 1)
		msg = strings.Replace(msg, "\u0001", "", 1)
		m.Text = msg
	}
}

func parseNewSub(m *common.Msg) {
	m.Type = common.MsgSub
	m.User.DisplayName = strings.Split(m.Text, " ")[0]
	m.User.Name = strings.ToLower(m.User.DisplayName)
}
