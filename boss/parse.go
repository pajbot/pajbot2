package boss

import (
	"log"
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
	}
	parseTags := true

	var splitline []string
	if strings.HasPrefix(line, ":") {
		parseTags = false
		splitline = strings.SplitN(line, ":", 2)
	} else {
		splitline = strings.SplitN(line, " :", 2)
	}
	tagsRaw := splitline[0]
	msg := splitline[1]
	tags := make(map[string]string)

	p.getMessage(msg)
	if p.m.User.Name == "twitchnotify" {
		if !strings.Contains(p.m.Message, " to ") && !strings.Contains(p.m.Message, " while ") {
			p.m.Type = common.MsgSub
			p.sub()
		} else {
			p.m.Type = common.MsgThrowAway
		}

	} else {
		tSplit := strings.Split(msg, " ")
		if len(tSplit) >= 2 {
			switch tSplit[1] {
			case "PRIVMSG":
				p.m.Type = common.MsgPrivmsg
				break
			case "WHISPER":
				p.m.Type = common.MsgWhisper
				break
			default:
				p.m.Type = common.MsgUnknown
				break
			}
		} else {
			p.m.Type = common.MsgUnknown
		}

		if p.m.Type == common.MsgUnknown {
			log.Printf("Unknown msg[%d]: %s", p.m.Type, msg)
		}

		// Should user properties stay at their zero value when there are no tags? Do we even care about this scenario?
		if parseTags {
			for _, tagValue := range strings.Split(tagsRaw, ";") {
				spl := strings.Split(tagValue, "=")
				k := spl[0]
				v := spl[1]
				tags[k] = v
			}
			p.getTwitchEmotes(tags["emotes"])
			p.getTags(tags)
		}
	}

	if p.m.Channel == p.m.User.Name {
		p.m.User.ChannelOwner = true
	}

	return *p.m
}

func (p *parse) getTwitchEmotes(emotetag string) {
	// TODO: Parse more emote information (bttv (and ffz?), name, size, isGif)
	// will we done by a module in the bot itself
	p.m.Emotes = make([]common.Emote, 0)
	if emotetag == "" {
		return
	}
	emoteSlice := strings.Split(emotetag, "/")
	for i := range emoteSlice {
		id := strings.Split(emoteSlice[i], ":")[0]
		e := &common.Emote{}
		e.Type = "twitch"
		e.Name = ""
		e.ID = id
		// 28 px should be fine for twitch emotes
		e.SizeX = 28
		e.SizeY = 28
		e.Count = strings.Count(emoteSlice[i], "-")
		p.m.Emotes = append(p.m.Emotes, *e)
	}
}

func (p *parse) getTags(tags map[string]string) {
	// TODO: Parse id and color
	// color and id is pretty useless imo
	if tags["display-name"] == "" {
		p.m.User.DisplayName = p.m.User.Name
	} else {
		p.m.User.DisplayName = tags["display-name"]
	}
	p.m.User.Type = tags["user-type"]
	if tags["turbo"] == "1" {
		p.m.User.Turbo = true
	}
	if tags["mod"] == "1" {
		p.m.User.Mod = true
	}
	if tags["subscriber"] == "1" {
		p.m.User.Sub = true
	}

}

func (p *parse) getMessage(msg string) {
	if strings.HasPrefix(msg, ":") {
		msg = strings.Replace(msg, ":", "", 1)
	}
	mSplit := strings.SplitN(msg, " :", 2)
	if len(mSplit) >= 2 {
		p.m.Message = strings.SplitN(msg, " :", 2)[1]
	}
	p.m.User.Name = strings.SplitN(msg, "!", 2)[0]
	cSplit := strings.SplitN(msg, "#", 3)
	if len(cSplit) >= 2 {
		p.m.Channel = strings.SplitN(cSplit[1], " ", 2)[0]
	}
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
