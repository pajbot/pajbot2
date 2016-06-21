package filter

import "github.com/pajlada/pajbot2/common"

// Filter is the shared interface for all filters
type Filter interface {
	Run(m string, msg *common.Msg, action *BanAction)
}
