package base

import (
	"sync"
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

type Command struct {
	cooldownMutex  *sync.RWMutex
	userCooldowns  map[string]bool
	globalCooldown bool

	UserCooldown   int
	GlobalCooldown int

	Description string
}

func New() Command {
	return Command{
		cooldownMutex: &sync.RWMutex{},
		userCooldowns: make(map[string]bool),

		UserCooldown:   15,
		GlobalCooldown: 5,
	}
}

func (c *Command) addCooldown(userID string) {
	c.cooldownMutex.Lock()
	c.userCooldowns[userID] = true
	c.globalCooldown = true
	c.cooldownMutex.Unlock()
}

func (c *Command) removeGlobalCooldown() {
	c.cooldownMutex.Lock()
	c.globalCooldown = false
	c.cooldownMutex.Unlock()
}

func (c *Command) removeCooldown(userID string) {
	c.cooldownMutex.Lock()
	delete(c.userCooldowns, userID)
	c.cooldownMutex.Unlock()
}

func (c *Command) hasCooldown(userID string) (ok bool) {
	c.cooldownMutex.RLock()
	if c.globalCooldown {
		ok = true
	} else {
		_, ok = c.userCooldowns[userID]
	}
	c.cooldownMutex.RUnlock()
	return ok
}

func (c *Command) HasCooldown(user pkg.User) bool {
	return c.hasCooldown(user.GetID())
}

func (c *Command) HasUserIDCooldown(userID string) bool {
	return c.hasCooldown(userID)
}

func (c *Command) AddCooldown(user pkg.User) {
	c.addCooldown(user.GetID())

	time.AfterFunc(time.Duration(c.UserCooldown)*time.Second, func() {
		c.removeCooldown(user.GetID())
	})
	time.AfterFunc(time.Duration(c.GlobalCooldown)*time.Second, func() {
		c.removeGlobalCooldown()
	})
}

func (c *Command) AddUserIDCooldown(userID string) {
	c.cooldownMutex.Lock()
	c.userCooldowns[userID] = true
	c.cooldownMutex.Unlock()

	time.AfterFunc(time.Duration(c.UserCooldown)*time.Second, func() {
		c.removeCooldown(userID)
	})
}

func (c *Command) AddGlobalCooldown() {
	c.cooldownMutex.Lock()
	c.globalCooldown = true
	c.cooldownMutex.Unlock()

	time.AfterFunc(time.Duration(c.GlobalCooldown)*time.Second, func() {
		c.removeGlobalCooldown()
	})
}
