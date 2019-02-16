package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string

	moderatorRoles = []string{
		"Snus Addict",
		"Roleplayer",
		"Moderators",
	}
)

const (
	// TODO: Make this a choice somewhere :pepega:
	moderationActionChannelID = `516960063081021460`
)

func init() {
	token = os.Getenv("PAJBOT2_DISCORD_BOT_TOKEN")

	if token == "" {
		fmt.Println("Missing bot token")
		os.Exit(1)
	}
}

func main() {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	bot.AddHandler(onMessage)

	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	defer bot.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// returns true if the given user id is in one of the given roles
func memberInRoles(s *discordgo.Session, guildID string, userID string, roles []string) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		for _, tRole := range roles {
			if role.Name == tRole {
				return true, nil
			}
		}
	}

	return false, nil
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	parts := strings.Split(m.Content, " ")

	if parts[0] == "$ban" {
		hasAccess, err := memberInRoles(s, m.GuildID, m.Author.ID, moderatorRoles)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if hasAccess {
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
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you don't have permission dummy", m.Author.Mention()))

		return
	}

	if parts[0] == "$channelinfo" {
		msg := fmt.Sprintf("Channel ID: %s", m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}

	// if not pajlada xd
	if m.Author.ID != "85699361769553920" {
		return
	}

	if parts[0] == "$guildinfo" {
		msg := fmt.Sprintf("Server ID: %s", m.GuildID)
		s.ChannelMessageSend(m.ChannelID, msg)
		return
	}
}
