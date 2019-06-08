package modules

import (
	"regexp"

	"github.com/pajbot/pajbot2/pkg"
	xurls "mvdan.cc/xurls/v2"
)

func init() {
	Register("link_filter", func() pkg.ModuleSpec {
		relaxedRegexp := xurls.Relaxed()
		strictRegexp := xurls.Strict()

		return &moduleSpec{
			id:   "link_filter",
			name: "Link filter",
			maker: func(b base) pkg.Module {
				return newLinkFilter(b, relaxedRegexp, strictRegexp)
			},
		}
	})
}

type LinkFilter struct {
	base

	relaxedRegexp *regexp.Regexp
	strictRegexp  *regexp.Regexp
}

func newLinkFilter(b base, relaxedRegexp, strictRegexp *regexp.Regexp) pkg.Module {
	return &LinkFilter{
		base: b,

		relaxedRegexp: relaxedRegexp,
		strictRegexp:  strictRegexp,
	}
}

func (m LinkFilter) OnMessage(bot pkg.BotChannel, source pkg.User, message pkg.Message, action pkg.Action) error {
	if source.IsModerator() || source.IsBroadcaster(bot.Channel()) {
		return nil
	}

	links := m.relaxedRegexp.FindAllString(message.GetText(), -1)
	if len(links) > 0 {
		action.Set(pkg.Timeout{180, "No links allowed"})
	}

	return nil
}
