package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	c2 "github.com/pajlada/pajbot2/pkg/commands"
)

var _ Command = &cmdGuildInfo{}

type cmdGuildInfo struct {
	c2.Base
}

func newGuildInfo() *cmdGuildInfo {
	return &cmdGuildInfo{
		Base: c2.NewBase(),
	}
}

func (c *cmdGuildInfo) Run(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) CommandResult {
	msg := fmt.Sprintf("Server ID: %s", m.GuildID)
	s.ChannelMessageSend(m.ChannelID, msg)
	return CommandResultFullCooldown
}
