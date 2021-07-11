package link_filter

import (
	"regexp"
	"time"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	xurls "mvdan.cc/xurls/v2"
)

func init() {
	modules.Register("link_filter", func() pkg.ModuleSpec {
		const regexpModifier = `(\b|$)`
		relaxedRegexpStr := xurls.Relaxed().String()
		strictRegexpStr := xurls.Strict().String()

		relaxedRegexp := regexp.MustCompile(relaxedRegexpStr + regexpModifier)
		relaxedRegexp.Longest()
		strictRegexp := regexp.MustCompile(strictRegexpStr + regexpModifier)
		strictRegexp.Longest()

		return modules.NewSpec("link_filter", "Link filter", false, func(b *mbase.Base) pkg.Module {
			return newLinkFilter(b, relaxedRegexp, strictRegexp)
		})
	})
}

type LinkFilter struct {
	mbase.Base

	relaxedRegexp *regexp.Regexp
	strictRegexp  *regexp.Regexp
}

func newLinkFilter(b *mbase.Base, relaxedRegexp, strictRegexp *regexp.Regexp) pkg.Module {
	return &LinkFilter{
		Base: *b,

		relaxedRegexp: relaxedRegexp,
		strictRegexp:  strictRegexp,
	}
}

func (m *LinkFilter) checkMessage(text string) bool {
	links := m.relaxedRegexp.FindAllString(text, -1)
	return len(links) > 0
}

func (m LinkFilter) OnMessage(event pkg.MessageEvent) pkg.Actions {
	if event.User.IsModerator() {
		return nil
	}

	if event.User.IsVIP() {
		return nil
	}

	if m.checkMessage(event.Message.GetText()) {
		actions := &twitchactions.Actions{}
		actions.Timeout(event.User, 180*time.Second).SetReason("No links allowed")
		return actions
	}

	return nil
}
