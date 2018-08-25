package modules

import (
	"github.com/pajlada/pajbot2/pkg"
	"mvdan.cc/xurls"
)

type LinkFilter struct {
}

func NewLinkFilter() *LinkFilter {
	return &LinkFilter{}
}

func (m *LinkFilter) Register() error {
	return nil
}

func (m LinkFilter) Name() string {
	return "LinkFilter"
}

func (m LinkFilter) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m LinkFilter) OnMessage(bot pkg.Sender, channel pkg.Channel, source pkg.User, message pkg.Message, action pkg.Action) error {
	if source.IsModerator() || source.IsBroadcaster(channel) {
		return nil
	}

	if channel.GetChannel() != "forsen" {
		return nil
	}

	links := xurls.Relaxed().FindAllString(message.GetText(), -1)
	if len(links) > 0 {
		action.Set(pkg.Timeout{180, "No links allowed"})
	}

	return nil
}
