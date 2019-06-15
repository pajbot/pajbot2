package modules

import (
	"regexp"
	"time"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
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

func (m LinkFilter) OnMessage(event pkg.MessageEvent) pkg.Actions {
	if event.User.IsModerator() {
		return nil
	}

	links := m.relaxedRegexp.FindAllString(event.Message.GetText(), -1)
	if len(links) > 0 {
		actions := &twitchactions.Actions{}
		actions.Timeout(event.User, 180*time.Second).SetReason("No links allowed")
		return actions
	}

	return nil
}
