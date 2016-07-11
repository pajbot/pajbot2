package command

import (
	"sync"
	"time"

	"github.com/pajlada/pajbot2/common"
)

// BaseCommand xD
type BaseCommand struct {
	sync.Mutex   // we probably need this right?
	ID           int
	Triggers     []string
	Level        int
	Cooldown     time.Duration
	UserCooldown time.Duration
	lastUse      map[string]time.Time // lastUse["global"] is for everyone
}

// OnCooldown checks if the command is on cooldown
func (c *BaseCommand) OnCooldown(user *common.User) bool {
	var usercd time.Duration
	var cd time.Duration
	switch {
	case user.Level >= 2000:
		cd, usercd = 0, 0
	case user.Level >= 1000:
		cd = c.Cooldown / 10
		usercd = c.UserCooldown / 10
	case user.Level >= 500:
		cd = c.Cooldown / 5
		usercd = c.UserCooldown / 5
	default:
		cd = c.Cooldown
		usercd = c.UserCooldown
	}
	now := time.Now()
	var onCD bool
	c.Lock()
	defer c.Unlock()
	if lastuse, ok := c.lastUse["global"]; ok {
		if time.Since(lastuse) < cd {
			return true
		}
	}

	if lastuse, ok := c.lastUse[user.Name]; ok {
		if time.Since(lastuse) < usercd {
			return true
		}
	}

	if !onCD {
		c.lastUse["global"] = now
		c.lastUse[user.Name] = now
	}
	return onCD
}
