package modules

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/pajlada/pajbot2/bot"

	"github.com/pajlada/pajbot2/common"
)

// Pyramid module from pajbot 1
type Pyramid struct {
	data      [][]string
	goingDown bool
}

// Ensure the module implements the interface properly
var _ Module = (*Pyramid)(nil)

// Check KKona
func (module *Pyramid) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if msg.User.Name == "twitchnotify" {
		return nil
	}
	msgParts := strings.Split(msg.Message, " ")
	if len(module.data) > 0 {
		curLen := len(msgParts)
		lastLen := len(module.data[len(module.data)-1])
		pyramidThing := module.data[0][0]
		lenDiff := curLen - lastLen
		log.Println(module.data)
		log.Println(msgParts)

		if math.Abs(float64(lenDiff)) == 1 {
			good := true

			for _, w := range msgParts {
				if w != pyramidThing {
					good = false
					break
				}
			}
			if good {
				module.data = append(module.data, msgParts)
				if lenDiff > 0 {
					if module.goingDown {
						module.data = make([][]string, 0)
						module.goingDown = false
					}
				} else if lenDiff < 0 {
					module.goingDown = true
					if curLen == 1 {
						// a pyramid has been finished
						var peakLen int
						for _, x := range module.data {
							if len(x) > peakLen {
								peakLen = len(x)
								log.Println(peakLen)
							}
						}
						if peakLen > 2 {
							m := fmt.Sprintf("%s just finished a %d width %s pyramid PogChamp pajaClap",
								msg.User.DisplayName,
								peakLen,
								pyramidThing)
							action.Response = m
						}
					}
				}
			}
		} else {
			module.data = make([][]string, 0)
			module.goingDown = false
		}
	}
	if len(msgParts) == 1 && len(module.data) == 0 {
		module.data = append(module.data, msgParts)
	}
	return nil
}
