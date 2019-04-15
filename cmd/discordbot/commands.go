package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	c2 "github.com/pajlada/pajbot2/pkg/commands"
)

// CommandResult xd
type CommandResult int

const (
	// CommandResultNoCooldown xd
	CommandResultNoCooldown CommandResult = iota
	// CommandResultUserCooldown xd
	CommandResultUserCooldown
	// CommandResultGlobalCooldown xd
	CommandResultGlobalCooldown
	// CommandResultFullCooldown xd
	CommandResultFullCooldown
)

// Command xd
type Command interface {
	HasUserIDCooldown(string) bool
	AddUserIDCooldown(string)
	AddGlobalCooldown()
	Run(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) CommandResult
}

var _ Command = &cmdPing{}

type cmdPing struct {
	c2.Base
}

func newPing() *cmdPing {
	return &cmdPing{
		Base: c2.NewBase(),
	}
}

func (c *cmdPing) Run(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) CommandResult {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, pong", m.Author.Mention()))
	return CommandResultFullCooldown
}
