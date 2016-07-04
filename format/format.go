package format

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/redismanager"
)

// Command is not a good name, but idk what else to call it
type Command struct {
	c       string
	subC    []string
	rawCmd  string
	outcome string
}

var mainRegex = regexp.MustCompile(`\$\([a-z\.]+\)`)
var partRegex = regexp.MustCompile(`[a-z]+`)

/*
ParseLine parses all variables (i.e. $(user)) from a line and returns
the new line, along with a list of commands
*/
func ParseLine(line string) (string, []Command) {
	log.Debug(line)
	matches := mainRegex.FindAllString(line, -1)
	log.Debug(matches)
	var cmds []Command
	for _, match := range matches {
		cmdlist := partRegex.FindAllString(match, -1)
		c := Command{
			c: cmdlist[0],
		}
		if len(cmdlist) > 1 {
			c.subC = cmdlist[1:]
		}
		c.rawCmd = match
		// lazy fix to avoid out of range error
		c.subC = append(c.subC, "")
		c.subC = append(c.subC, "")
		cmds = append(cmds, c)
	}
	log.Debug(cmds, line)
	return line, cmds
}

/*
RunCommands applies a list of commands on the given line
*/
func RunCommands(line string, cmds []Command) string {
	for _, c := range cmds {
		line = strings.Replace(line, c.rawCmd, c.outcome, 1)
	}
	return line
}

/*
ExecCommand xD
*/
func ExecCommand(redis *redismanager.RedisManager, cmd *Command, msg *common.Msg) {
	switch cmd.c {
	case "source", "sender":
		cmd.outcome = ParseUser(&msg.User, cmd.subC)
	case "user":
		if msg.Args != nil {
			if redis.IsValidUser(msg.Channel, msg.Args[0]) {
				user := redis.LoadUser(msg.Channel, msg.Args[0])
				cmd.outcome = ParseUser(&user, cmd.subC)
				return
			}
		}
		cmd.outcome = ParseUser(&msg.User, cmd.subC)
	}
}

/*
ParseUser xD
*/
func ParseUser(user *common.User, cmds []string) string {
	if cmds == nil {
		return user.DisplayName
	}
	switch cmds[0] {
	case "name":
		return user.Name
	case "points":
		return fmt.Sprintf("%d", user.Points)
	case "level":
		return fmt.Sprintf("%d", user.Level)
	case "lines":
		return ParseLines(user, cmds[1])
	default:
		return user.DisplayName
	}
}

/*
ParseLines xD
*/
func ParseLines(user *common.User, arg string) string {
	switch arg {
	case "online":
		return fmt.Sprintf("%d", user.OnlineMessageCount)
	case "offline":
		return fmt.Sprintf("%d", user.OfflineMessageCount)
	default:
		return fmt.Sprintf("%d", user.TotalMessageCount)
	}
}
