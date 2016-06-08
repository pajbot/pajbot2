package boss

import (
	"fmt"
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
func Parse(line string) common.Msg {
	p := &parse{}
	p.m = &common.Msg{
		User: common.User{},
	}
	parseTags := true
	fmt.Println(line)
	if strings.Contains(line, "twitchnotify") {
		fmt.Println(line)
	}
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

	p.GetMessage(msg)
	if p.m.User.Name == "twitchnotify" {
		p.m.Type = "sub"
		p.Sub()
	} else {
		if strings.Contains(msg, "PRIVMSG") {
			p.m.Type = "privmsg"
		} else {
			p.m.Type = "whisper"
		}

		if parseTags {
			for _, tagValue := range strings.Split(tagsRaw, ";") {
				spl := strings.Split(tagValue, "=")
				k := spl[0]
				v := spl[1]
				tags[k] = v
			}
			p.GetTwitchEmotes(tags["emotes"])
			p.GetTags(tags)
		}
	}

	return *p.m
}

func (p *parse) GetTwitchEmotes(emotetag string) {
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
		e.Count = strings.Count(emoteSlice[i], "-")
		p.m.Emotes = append(p.m.Emotes, *e)
	}
}

func (p *parse) GetTags(tags map[string]string) {
	p.m.User.Displayname = tags["display-name"]
	p.m.User.Type = tags["user-type"]

	if tags["turbo"] == "1" {
		p.m.User.Turbo = true
	}
	if tags["mod"] == "1" {
		p.m.User.Mod = true
	}
	if tags["subscriber"] == "1" {
		p.m.User.Mod = true
	}

}

func (p *parse) GetMessage(msg string) {
	if strings.HasPrefix(msg, ":") {
		msg = strings.Replace(msg, ":", "", 1)
	}
	//fmt.Println(msg)
	p.m.Message = strings.SplitN(msg, " :", 2)[1]
	p.m.User.Name = strings.SplitN(msg, "!", 2)[0]
	c := strings.SplitN(msg, "#", 3)[1]
	p.m.Channel = strings.SplitN(c, " ", 2)[0]
}

func (p *parse) Sub() {
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
	p.m.User.Displayname = strings.Split(m, " ")[0]
	p.m.User.Name = strings.ToLower(p.m.User.Displayname)
}
