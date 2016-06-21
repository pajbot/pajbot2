package filter

import (
	"fmt"

	"github.com/mvdan/xurls"
	"github.com/pajlada/pajbot2/common"
)

// Link xD
type Link struct {
}

var _ Filter = (*Link)(nil)

// Run filter
func (Link *Link) Run(m string, msg *common.Msg, action *BanAction) {
	matches := xurls.Relaxed.FindAllString(m, -1)
	action.Matches = matches
	if len(matches) > 0 {
		action.Level = 3
		action.Reason = fmt.Sprintf("matched links: %v", matches)
		action.Matched = true
	}
}
