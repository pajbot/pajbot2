package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pajlada/pajbot2/internal/config"
	"github.com/pajlada/pajbot2/pkg/commandmatcher"
	"github.com/pajlada/stupidmigration"

	_ "github.com/lib/pq"
)

var (
	token string

	adminRole         = "104258168128802816"
	moderatorRole     = "132930561361707010"
	miniModeratorRole = "531783703882629126"
	mutedRole         = "431734043881504778"

	adminRoles = []string{
		adminRole,
	}

	moderatorRoles = []string{
		adminRole,
		moderatorRole,
	}

	miniModeratorRoles = []string{
		adminRole,
		moderatorRole,
		miniModeratorRole,
	}

	commands = commandmatcher.NewMatcher()
)

const (
	// TODO: Make this a choice somewhere :pepega:
	moderationActionChannelID = `516960063081021460`

	weebChannelID = `500650560614301696`
)

var sqlClient *sql.DB

func init() {
	token = os.Getenv("PAJBOT2_DISCORD_BOT_TOKEN")

	if token == "" {
		fmt.Println("Missing bot token")
		os.Exit(1)
	}

	var err error
	sqlClient, err = sql.Open("postgres", config.GetDSN())
	if err != nil {
		fmt.Println("Unable to connect to mysql", err)
		os.Exit(1)
	}

	err = sqlClient.Ping()
	if err != nil {
		fmt.Println("Unable to ping mysql", err)
		os.Exit(1)
	}

	err = stupidmigration.Migrate("migrations", sqlClient)
	if err != nil {
		fmt.Println("Unable to run SQL migrations", err)
		os.Exit(1)
	}
}

func registerCommands() {
	// TODO: unmute

	commands.Register([]string{"$mute"}, newMute())

	commands.Register([]string{"$ping"}, newPing())

	commands.Register([]string{"$ban"}, newBan())

	commands.Register([]string{"$channelinfo"}, newChannelInfo())
	commands.Register([]string{"$guildinfo", "$serverinfo"}, newGuildInfo())

	commands.Register([]string{"$test-minimod"}, func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) {
		hasAccess, err := memberInRoles(s, m.GuildID, m.Author.ID, miniModeratorRoles)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !hasAccess {
			// No permission
			return
		}

		s.ChannelMessageSend(m.ChannelID, "you are >= minimod")
	})

	commands.Register([]string{"$test-mod"}, func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) {
		hasAccess, err := memberInRoles(s, m.GuildID, m.Author.ID, moderatorRoles)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !hasAccess {
			// No permission
			return
		}

		s.ChannelMessageSend(m.ChannelID, "you are >= mod")
	})

	commands.Register([]string{"$test-admin"}, func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) {
		hasAccess, err := memberInRoles(s, m.GuildID, m.Author.ID, adminRoles)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !hasAccess {
			// No permission
			return
		}

		s.ChannelMessageSend(m.ChannelID, "you are >= admin")
	})

	commands.Register([]string{"$roles"}, func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) {
		if m.Author.ID != "85699361769553920" {
			return
		}

		roles, err := s.GuildRoles(m.GuildID)
		if err != nil {
			fmt.Println("Error getting roles:", err)
			return
		}

		response := "```"
		for _, role := range roles {
			if role.Managed {
				continue
			}
			response += fmt.Sprintf("%s = %s\n", role.ID, role.Name)
		}
		response += "```"

		s.ChannelMessageSend(m.ChannelID, response)
	})

	commands.Register([]string{"$userid"}, func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string) {
		hasAccess, err := memberInRoles(s, m.GuildID, m.Author.ID, miniModeratorRoles)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !hasAccess {
			// No permission
			return
		}

		if len(parts) < 2 {
			return
		}

		target := parts[1]
		targetID := cleanUserID(parts[1])

		s.ChannelMessageSend(m.ChannelID, "User ID of "+target+" is `"+targetID+"`")
	})
}

func main() {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	go func() {
		for {
			<-time.After(3 * time.Second)
			now := time.Now()
			const query = `SELECT id, action, timepoint FROM discord_queue ORDER BY timepoint DESC LIMIT 30;`
			rows, err := sqlClient.Query(query)
			if err != nil {
				fmt.Println("err:", err)
				continue
			}
			var actionsToRemove []int64
			defer rows.Close()
			for rows.Next() {
				var (
					id           int64
					actionString string
					timepoint    time.Time
				)
				if err := rows.Scan(&id, &actionString, &timepoint); err != nil {
					fmt.Println("Error scanning:", err)
				}
				if timepoint.After(now) {
					continue
				}

				var action Action
				err = json.Unmarshal([]byte(actionString), &action)
				if err != nil {
					fmt.Println("Error unmarshaling action:", err)
					continue
				}
				fmt.Println("Perform", action.Type)

				switch action.Type {
				case "unmute":
					err = bot.GuildMemberRoleRemove(action.GuildID, action.UserID, action.RoleID)
					if err != nil {
						fmt.Println("Error removing role")
						continue
					}

					actionsToRemove = append(actionsToRemove, id)
				}
			}

			for _, actionID := range actionsToRemove {
				sqlClient.Exec("DELETE FROM discord_queue WHERE id=$1", actionID)
			}
		}
	}()

	registerCommands()

	bot.AddHandler(onMessage)
	bot.AddHandler(onUserBanned)
	bot.AddHandler(onMessageReactionAdded)
	bot.AddHandler(onMessageReactionRemoved)

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
			if role.ID == tRole {
				return true, nil
			}
		}
	}

	return false, nil
}

var patternUserIDReplacer = regexp.MustCompile(`^<@!?([0-9]+)>$`)
var patternUserID = regexp.MustCompile(`^[0-9]+$`)

func cleanUserID(input string) string {
	output := patternUserIDReplacer.ReplaceAllString(input, "$1")

	if !patternUserID.MatchString(output) {
		return ""
	}

	return output
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	c, parts := commands.Match(m.Content)
	if c != nil {
		if cmd, ok := c.(Command); ok {
			id := m.ChannelID + m.Author.ID
			if cmd.HasUserIDCooldown(id) {
				return
			}

			switch cmd.Run(s, m, parts) {
			case CommandResultUserCooldown:
				cmd.AddUserIDCooldown(id)
			case CommandResultGlobalCooldown:
				cmd.AddGlobalCooldown()
			case CommandResultFullCooldown:
				cmd.AddUserIDCooldown(id)
				cmd.AddGlobalCooldown()
			}
		} else if f, ok := c.(func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string)); ok {
			f(s, m, parts)
		}
	}
}

func onUserBanned(s *discordgo.Session, m *discordgo.GuildBanAdd) {
	auditLog, err := s.GuildAuditLog(m.GuildID, "", "", 22, 1)
	if err != nil {
		fmt.Println("Error getting user ban data", err)
		return
	}
	fmt.Println(auditLog)
	if len(auditLog.AuditLogEntries) != 1 {
		fmt.Println("Unable to get the single ban entry")
		return
	}
	if len(auditLog.Users) != 2 {
		fmt.Println("length of users is wrong")
		return
	}
	banner := auditLog.Users[0]
	bannedUser := auditLog.Users[1]
	if bannedUser.ID != m.User.ID {
		fmt.Println("got log for wrong use Pepega")
		return
	}
	fmt.Println(auditLog.Users)
	entry := auditLog.AuditLogEntries[0]
	// var username string
	// for _ user := range auditLog.Users {
	// 	if user.ID == entry.
	// }
	fmt.Println(entry)
	fmt.Println("Entry User ID:", entry.UserID)
	fmt.Println("target user ID:", m.User.ID)
	s.ChannelMessageSend(moderationActionChannelID, fmt.Sprintf("%s was banned by %s: %s", m.User.Mention(), banner.Username, entry.Reason))
}

// const weebMessageID = `552788256333234176`
const weebMessageID = `552791672854151190`
const reactionBye = "👋"

func onMessageReactionAdded(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.MessageID == weebMessageID {
		if m.Emoji.Name == reactionBye {
			c, err := s.State.Channel(weebChannelID)
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			var overwriteDenies int
			for _, overwrite := range c.PermissionOverwrites {
				if overwrite.Type == "member" && overwrite.ID == m.UserID {
					overwriteDenies = overwrite.Deny
				}
			}
			if overwriteDenies != 0 {
				// s.ChannelMessageSend(m.ChannelID, "cannot set your permissions - you have weird permissions set from before")
				return
			}

			err = s.ChannelPermissionSet(weebChannelID, m.UserID, "member", 0, discordgo.PermissionReadMessages)
			if err != nil {
				fmt.Println("uh oh something went wrong")
				return
			}

			// s.ChannelMessageSend(m.ChannelID, "added permission")
		}
	}
}

func onMessageReactionRemoved(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m.MessageID == weebMessageID {
		if m.Emoji.Name == reactionBye {
			c, err := s.State.Channel(weebChannelID)
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			var overwriteDenies int
			for _, overwrite := range c.PermissionOverwrites {
				if overwrite.Type == "member" && overwrite.ID == m.UserID {
					overwriteDenies = overwrite.Deny
				}
			}

			if overwriteDenies != discordgo.PermissionReadMessages {
				// s.ChannelMessageSend(m.ChannelID, "not allowed to remove that permission buddy")
				return
			}

			err = s.ChannelPermissionDelete(weebChannelID, m.UserID)
			if err != nil {
				fmt.Println("uh oh something went wrong")
				return
			}
			// s.ChannelMessageSend(m.ChannelID, "removed permission")
		}
	}
}
