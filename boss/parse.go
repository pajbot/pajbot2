package boss

import (
	"strconv"
	"strings"

	"github.com/pajlada/pajbot2/common"
)

type parse struct {
	m *common.Msg
}

/*
Parse parses an IRC message into a more readable bot.Msg
*/
func (p *parse) Parse(line string) common.Msg {
	p.m = &common.Msg{
		User: common.User{},
		Type: common.MsgUnknown,
	}

	// msg is the string we will keep working on/reducing as we parse things
	msg := line

	var splitLine []string

	// The message starts with @, that means there are IRCv3 tags available to parse
	if strings.HasPrefix(line, "@") {
		splitLine = strings.SplitN(msg, " ", 2)
		p.parseTags(splitLine[0])
		msg = splitLine[1]
	}

	// Parse source
	splitLine = strings.SplitN(msg, " ", 2)
	p.parseSource(splitLine[0])
	msg = splitLine[1]

	// Parse message type
	splitLine = strings.SplitN(msg, " ", 2)
	p.parseMsgType(splitLine[0])
	msg = splitLine[1]

	if p.m.Type == common.MsgUnknown {
		p.m.Type = common.MsgThrowAway
		return *p.m
	}

	splitLine = strings.SplitN(msg, " ", 2)
	p.parseChannel(splitLine[0])
	msg = splitLine[1]

	// Parse message + msg type (if it's a /me message or not)
	p.parseMessage(msg)

	// TODO: fix this sub detection (@pajlada)
	if p.m.User.Name == "twitchnotify" {
		if !strings.Contains(p.m.Message, " to ") && !strings.Contains(p.m.Message, " while ") {
			p.m.Type = common.MsgSub
			p.sub()
		}
	}

	// If the destination of the message is the same as the username,
	// then we tag the user as the channel owner. This will automatically
	// give him access to broadcaster commands
	if p.m.Channel == p.m.User.Name {
		p.m.User.ChannelOwner = true
	}

	if p.m.Tags != nil {
		// Parse tags further, such as the msg-id value for determinig the msg type
		p.parseExtendedTags()
	}

	return *p.m
}

func (p *parse) parseTwitchEmotes(emotetag string) {
	// TODO: Parse more emote information (bttv (and ffz?), name, size, isGif)
	// will we done by a module in the bot itself
	p.m.Emotes = make([]common.Emote, 0)
	if emotetag == "" {
		return
	}
	emoteSlice := strings.Split(emotetag, "/")
	for i := range emoteSlice {
		spl := strings.Split(emoteSlice[i], ":")
		id := spl[0]
		e := &common.Emote{}
		e.Type = "twitch"
		e.Name = p.getEmoteName(spl[1])
		e.ID = id
		// 28 px should be fine for twitch emotes
		e.SizeX = 28
		e.SizeY = 28
		e.Count = strings.Count(emoteSlice[i], "-")
		p.m.Emotes = append(p.m.Emotes, *e)
	}
}

func (p *parse) getEmoteName(pos string) string {
	pos = strings.Split(pos, ",")[0]
	spl := strings.Split(pos, "-")
	start, _ := strconv.Atoi(spl[0])
	end, _ := strconv.Atoi(spl[1])
	runes := []rune(p.m.Message)
	name := runes[start : end+1]
	return string(name)
}

func (p *parse) parseTagValues() {
	// TODO: Parse id and color
	// color and id is pretty useless imo
	if p.m.Tags["display-name"] == "" {
		p.m.User.DisplayName = p.m.User.Name
	} else {
		p.m.User.DisplayName = p.m.Tags["display-name"]
	}
	delete(p.m.Tags, "display-name")
	p.m.User.Type = p.m.Tags["user-type"]
	delete(p.m.Tags, "user-type")
	if p.m.Tags["turbo"] == "1" {
		p.m.User.Turbo = true
	}
	delete(p.m.Tags, "turbo")
	if p.m.Tags["mod"] == "1" {
		p.m.User.Mod = true
	}
	delete(p.m.Tags, "mod")
	if p.m.Tags["subscriber"] == "1" {
		p.m.User.Sub = true
	}
	delete(p.m.Tags, "subscriber")
}

func (p *parse) parseExtendedTags() {
	// Parse twitch emotes from the "emotes" tag
	p.parseTwitchEmotes(p.m.Tags["emotes"])
	delete(p.m.Tags, "emotes")

	switch p.m.Tags["msg-id"] {
	case "resub":
		p.m.Type = common.MsgReSub

	case "":
		break

	default:
		p.m.Type = common.MsgUnknown
	}

	if p.m.Tags["login"] != "" {
		p.m.User.Name = p.m.Tags["login"]
	}
}

/*
XXX: Should user properties stay at their zero value when there are no tags? Do we even care about this scenario?
*/
func (p *parse) parseTags(msg string) {
	p.m.Tags = make(map[string]string)
	// IRCv3-tags are separated by semicolons
	for _, tagValue := range strings.Split(msg, ";") {
		spl := strings.Split(tagValue, "=")
		k := spl[0]
		v := strings.Replace(spl[1], "\\s", " ", -1)
		p.m.Tags[k] = v
	}

	p.parseTagValues()

}

func (p *parse) parseSource(msg string) {
	msg = msg[1:]
	// Check if the source is a user
	userSepPos := strings.Index(msg, "!")
	hostSepPos := strings.Index(msg, "@")
	if userSepPos > -1 && hostSepPos > -1 && userSepPos < hostSepPos {
		// A valid user address is found!
		p.m.User.Name = msg[0:userSepPos]
	}
	log.Debug(msg)
}

func (p *parse) parseMsgType(msg string) {
	switch msg {
	case "PRIVMSG":
		p.m.Type = common.MsgPrivmsg

	case "WHISPER":
		p.m.Type = common.MsgWhisper

	case "USERNOTICE":
		p.m.Type = common.MsgUsernotice
	}
}

func (p *parse) parseChannel(msg string) {
	p.m.Channel = strings.Replace(msg[1:], "#", "", 0)
}

func (p *parse) parseMessage(msg string) {
	p.m.Message = msg[1:]

	// figure out whether the message is an ACTION or not
	p.getAction()
}

// regex in 2016 LUL
func (p *parse) getAction() {
	if strings.HasPrefix(p.m.Message, "\u0001ACTION ") && strings.HasSuffix(p.m.Message, "\u0001") {
		p.m.Me = true
		m := p.m.Message
		m = strings.Replace(m, "\u0001ACTION ", "", 1)
		m = strings.Replace(m, "\u0001", "", 1)
		p.m.Message = m
	}
}

// TODO: rewrite (@pajlada)
func (p *parse) sub() {
	m := p.m.Message
	if strings.Contains(m, "just ") {
		p.m.Length = 1
	} else {
		temp := strings.Split(m, " for ")[1]
		l := strings.Split(temp, " ")[0]
		length, err := strconv.Atoi(l)
		if err == nil {
			p.m.Length = length
		} else {
			panic(err)
		}
	}
	p.m.User.DisplayName = strings.Split(m, " ")[0]
	p.m.User.Name = strings.ToLower(p.m.User.DisplayName)
}
