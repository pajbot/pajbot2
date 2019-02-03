package modules

import (
	"github.com/pajlada/pajbot2/pkg"
	"mvdan.cc/xurls"
)

type LinkFilter struct {
	botChannel pkg.BotChannel
}

func newLinkFilter() pkg.Module {
	return &LinkFilter{}
}

var linkFilterSpec = moduleSpec{
	id:    "link_filter",
	name:  "Link filter",
	maker: newLinkFilter,
}

func (m *LinkFilter) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	return nil
}

func (m *LinkFilter) Disable() error {
	return nil
}

func (m *LinkFilter) Spec() pkg.ModuleSpec {
	return &linkFilterSpec
}

func (m *LinkFilter) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m LinkFilter) Name() string {
	return "LinkFilter"
}

func (m LinkFilter) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m LinkFilter) OnMessage(bot pkg.BotChannel, source pkg.User, message pkg.Message, action pkg.Action) error {
	if source.IsModerator() || source.IsBroadcaster(bot.Channel()) {
		return nil
	}

	links := xurls.Relaxed().FindAllString(message.GetText(), -1)
	if len(links) > 0 {
		action.Set(pkg.Timeout{180, "No links allowed"})
	}

	return nil
}
