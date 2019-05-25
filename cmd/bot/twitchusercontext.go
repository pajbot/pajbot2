package main

import (
	"fmt"
	"sync"

	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.UserContext = &UserContext{}

type UserContext struct {
	mutex *sync.Mutex

	// key = channel ID
	context map[string]map[string][]string
}

func NewUserContext() *UserContext {
	c := &UserContext{
		mutex:   &sync.Mutex{},
		context: make(map[string]map[string][]string),
	}

	return c
}

func (c *UserContext) GetContext(channelID, userID string) []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if users, ok := c.context[channelID]; ok {
		if userContext, ok := users[userID]; ok {
			return userContext
		}
	}

	return nil
}

func (c *UserContext) AddContext(channelID, userID, message string) {
	if channelID == "" {
		fmt.Println("Channel ID is empty")
		return
	}
	if userID == "" {
		fmt.Println("User ID is empty")
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.context[channelID]
	if !ok {
		c.context[channelID] = make(map[string][]string)
	}

	c.context[channelID][userID] = append(c.context[channelID][userID], message)
	newLen := len(c.context[channelID][userID]) - 5
	if newLen > 5 {
		c.context[channelID][userID] = c.context[channelID][userID][newLen:]
	}
}
