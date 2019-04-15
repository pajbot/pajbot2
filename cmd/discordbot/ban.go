package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	c2 "github.com/pajlada/pajbot2/pkg/commands"
)

var _ Command = &cmdBan{}

type cmdBan struct {
	c2.Base
}

func newBan() *cmdBan {
	return &cmdBan{
		Base: c2.NewBase(),
	}
}

func (c *cmdBan) Run(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) (res CommandResult) {
	res = CommandResultNoCooldown
	hasAccess, err := memberInRoles(s, m.GuildID, m.Author.ID, moderatorRoles)
	if err != nil {
		fmt.Println("Error:", err)
		return CommandResultUserCooldown
	}

	if !hasAccess {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you don't have permission dummy", m.Author.Mention()))
		return CommandResultUserCooldown
	}

	if len(m.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "missing user arg. usage: $ban <user> <reason>")
		return
	}

	target := m.Mentions[0]

	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "missing reason arg. usage: $ban <user> <reason>")
		return
	}

	reason := strings.Join(parts[2:], " ")

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Banning %s (%s) for reason: `%s`", target.Username, target.ID, reason))
	s.ChannelMessageSend(moderationActionChannelID, fmt.Sprintf("%s banned %s (%s) for reason: `%s`", m.Author.Username, target.Username, target.ID, reason))
	s.GuildBanCreateWithReason(m.GuildID, target.ID, reason, 0)

	return
}
