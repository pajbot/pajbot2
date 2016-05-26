package irc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nuuls/pajbot2/bot"
)

type parse struct {
	m *bot.Msg
}

func Parse(line string) bot.Msg {
	p := &parse{}
	p.m = &bot.Msg{}
	//fmt.Println(line)
	if strings.Contains(line, "twitchnotify") {
		fmt.Println(line)
	}
	var splitline []string
	if strings.HasPrefix(line, ":") {
		splitline = strings.SplitN(line, ":", 2)
	} else {
		splitline = strings.SplitN(line, " :", 2)
	}
	tagsRaw := splitline[0]
	msg := splitline[1]
	tags := make(map[string]string)

	p.GetMessage(msg)
	if p.m.Username == "twitchnotify" {
		p.m.MessageType = "sub"
		p.Sub()
	} else {
		if strings.Contains(msg, "PRIVMSG") {
			p.m.MessageType = "privmsg"
		} else {
			p.m.MessageType = "whisper"
		}
		tags_ := strings.Split(tagsRaw, ";")

		for i := range tags_ {
			k := strings.Split(tags_[i], "=")[0]
			v := strings.Split(tags_[i], "=")[1]
			tags[k] = v
		}
		p.GetTwitchEmotes(tags["emotes"])
		p.GetTags(tags)
	}

	return *p.m
}

func (p *parse) GetTwitchEmotes(emotetag string) {
	p.m.Emotes = make([]bot.Emote, 0)
	if emotetag == "" {
		return
	}
	emoteSlice := strings.Split(emotetag, "/")
	for i := range emoteSlice {
		id := strings.Split(emoteSlice[i], ":")[0]
		e := &bot.Emote{}
		e.EmoteType = "twitch"
		e.Name = ""
		e.Id = id
		e.Count = strings.Count(emoteSlice[i], "-")
		p.m.Emotes = append(p.m.Emotes, *e)
	}
}

func (p *parse) GetTags(tags map[string]string) {
	p.m.Color = tags["color"]
	p.m.Displayname = tags["display-name"]
	p.m.Usertype = tags["user-type"]

	if tags["turbo"] == "1" {
		p.m.Turbo = true
	}
	if tags["mod"] == "1" {
		p.m.Mod = true
	}
	if tags["subscriber"] == "1" {
		p.m.Subscriber = true
	}

}

func (p *parse) GetMessage(msg string) {
	if strings.HasPrefix(msg, ":") {
		msg = strings.Replace(msg, ":", "", 1)
	}
	//fmt.Println(msg)
	p.m.Message = strings.SplitN(msg, " :", 2)[1]
	p.m.Username = strings.SplitN(msg, "!", 2)[0]
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
	p.m.Displayname = strings.Split(m, " ")[0]
	p.m.Username = strings.ToLower(p.m.Displayname)
}
