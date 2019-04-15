package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	c2 "github.com/pajlada/pajbot2/pkg/commands"
)

var _ Command = &cmdChannelInfo{}

type cmdChannelInfo struct {
	c2.Base
}

func newChannelInfo() *cmdChannelInfo {
	return &cmdChannelInfo{
		Base: c2.NewBase(),
	}
}

func (c *cmdChannelInfo) Run(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) CommandResult {
	msg := fmt.Sprintf("Channel ID: %s", m.ChannelID)
	s.ChannelMessageSend(m.ChannelID, msg)
	return CommandResultFullCooldown
}
