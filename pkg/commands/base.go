package commands

import (
	"sync"
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

type Base struct {
	cooldownMutex  *sync.RWMutex
	userCooldowns  map[string]bool
	globalCooldown bool

	UserCooldown   int
	GlobalCooldown int

	Description string
}

func NewBase() Base {
	return Base{
		cooldownMutex: &sync.RWMutex{},
		userCooldowns: make(map[string]bool),

		UserCooldown:   15,
		GlobalCooldown: 5,
	}
}

func (c *Base) addCooldown(userID string) {
	c.cooldownMutex.Lock()
	c.userCooldowns[userID] = true
	c.globalCooldown = true
	c.cooldownMutex.Unlock()
}

func (c *Base) removeGlobalCooldown() {
	c.cooldownMutex.Lock()
	c.globalCooldown = false
	c.cooldownMutex.Unlock()
}

func (c *Base) removeCooldown(userID string) {
	c.cooldownMutex.Lock()
	delete(c.userCooldowns, userID)
	c.cooldownMutex.Unlock()
}

func (c *Base) hasCooldown(userID string) (ok bool) {
	c.cooldownMutex.RLock()
	if c.globalCooldown {
		ok = true
	} else {
		_, ok = c.userCooldowns[userID]
	}
	c.cooldownMutex.RUnlock()
	return ok
}

func (c *Base) HasCooldown(user pkg.User) bool {
	return c.hasCooldown(user.GetID())
}

func (c *Base) HasUserIDCooldown(userID string) bool {
	return c.hasCooldown(userID)
}

func (c *Base) AddCooldown(user pkg.User) {
	c.addCooldown(user.GetID())

	time.AfterFunc(time.Duration(c.UserCooldown)*time.Second, func() {
		c.removeCooldown(user.GetID())
	})
	time.AfterFunc(time.Duration(c.GlobalCooldown)*time.Second, func() {
		c.removeGlobalCooldown()
	})
}

func (c *Base) AddUserIDCooldown(userID string) {
	c.cooldownMutex.Lock()
	c.userCooldowns[userID] = true
	c.cooldownMutex.Unlock()

	time.AfterFunc(time.Duration(c.UserCooldown)*time.Second, func() {
		c.removeCooldown(userID)
	})
}

func (c *Base) AddGlobalCooldown() {
	c.cooldownMutex.Lock()
	c.globalCooldown = true
	c.cooldownMutex.Unlock()

	time.AfterFunc(time.Duration(c.GlobalCooldown)*time.Second, func() {
		c.removeGlobalCooldown()
	})
}
