package commands

import (
	"sync"
	"time"

	"github.com/pajlada/pajbot2/pkg"
)

type base struct {
	cooldownMutex  *sync.RWMutex
	userCooldowns  map[string]bool
	globalCooldown bool

	UserCooldown   int
	GlobalCooldown int
}

func newBase() base {
	return base{
		cooldownMutex: &sync.RWMutex{},
		userCooldowns: make(map[string]bool),

		UserCooldown:   15,
		GlobalCooldown: 5,
	}
}

func (c *base) addCooldown(userID string) {
	c.cooldownMutex.Lock()
	c.userCooldowns[userID] = true
	c.globalCooldown = true
	c.cooldownMutex.Unlock()
}

func (c *base) removeGlobalCooldown() {
	c.cooldownMutex.Lock()
	c.globalCooldown = false
	c.cooldownMutex.Unlock()
}

func (c *base) removeCooldown(userID string) {
	c.cooldownMutex.Lock()
	delete(c.userCooldowns, userID)
	c.cooldownMutex.Unlock()
}

func (c *base) hasCooldown(userID string) (ok bool) {
	c.cooldownMutex.RLock()
	if c.globalCooldown {
		ok = true
	} else {
		_, ok = c.userCooldowns[userID]
	}
	c.cooldownMutex.RUnlock()
	return ok
}

func (c *base) HasCooldown(user pkg.User) bool {
	return c.hasCooldown(user.GetID())
}

func (c *base) AddCooldown(user pkg.User) {
	c.addCooldown(user.GetID())

	time.AfterFunc(time.Duration(c.UserCooldown)*time.Second, func() {
		c.removeCooldown(user.GetID())
	})
	time.AfterFunc(time.Duration(c.GlobalCooldown)*time.Second, func() {
		c.removeGlobalCooldown()
	})
}
