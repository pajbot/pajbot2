package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	c2 "github.com/pajlada/pajbot2/pkg/commands"
	"github.com/pajlada/pajbot2/pkg/utils"
)

var _ Command = &cmdMute{}

type cmdMute struct {
	c2.Base
}

func newMute() *cmdMute {
	return &cmdMute{
		Base: c2.NewBase(),
	}
}

func (c *cmdMute) Run(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) (res CommandResult) {
	res = CommandResultNoCooldown
	const usage = `$mute @user duration <reason> (i.e. $mute @user 1h5m shitposting in serious channel)`

	var err error
	var targetID string
	var duration time.Duration
	var reason string

	hasAccess, err := memberInRoles(s, m.GuildID, m.Author.ID, miniModeratorRoles)
	if err != nil {
		fmt.Println("Error:", err)
		return CommandResultUserCooldown
	}
	if !hasAccess {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you don't have permission to use $mute", m.Author.Mention()))
		return CommandResultUserCooldown
	}

	parts = parts[1:]

	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" usage: "+usage)
		return
	}

	targetID = cleanUserID(parts[0])

	if targetID == "" {
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" $mute invalid user. usage: "+usage)
		return
	}

	duration, err = utils.ParseDuration(parts[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" $mute invalid duration: "+err.Error())
		return
	}

	if duration < 1*time.Minute {
		duration = 1 * time.Minute
	} else if duration > 14*24*time.Hour {
		duration = 14 * 24 * time.Hour
	}

	reason = strings.Join(parts[2:], " ")

	// Create queued up unmute action in database
	timepoint := time.Now().Add(duration)

	action := Action{
		Type:    "unmute",
		GuildID: m.GuildID,
		UserID:  targetID,
		RoleID:  mutedRole,
	}
	bytes, err := json.Marshal(&action)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" $mute unable to marshal action: "+err.Error())
		return
	}

	query := "INSERT INTO discord_queue (action, timepoint) VALUES ($1, $2)"
	_, err = sqlClient.Exec(query, string(bytes), timepoint)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" $mute sql error: "+err.Error())
		return
	}

	// Assign muted role
	err = s.GuildMemberRoleAdd(m.GuildID, targetID, mutedRole)
	if err != nil {
		fmt.Println("Error assigning role:", err)
	}

	// Announce mute in action channel
	s.ChannelMessageSend(moderationActionChannelID, fmt.Sprintf("%s muted %s for %s. reason: %s", m.Author.Mention(), targetID, duration, reason))
	fmt.Println(moderationActionChannelID, fmt.Sprintf("%s muted %s for %s. reason: %s", m.Author.Mention(), targetID, duration, reason))

	// Announce mute success
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s mute %s for %s. reason: %s", m.Author.Mention(), targetID, duration, reason))
	fmt.Println(m.ChannelID, fmt.Sprintf("%s mute %s for %s. reason: %s", m.Author.Mention(), targetID, duration, reason))

	return
}
