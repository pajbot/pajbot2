package format

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pajlada/pajbot2/common"
)

// Command is not a good name, but idk what else to call it
type Command struct {
	C       string
	SubC    []string
	RawCmd  string
	Outcome string
}

var mainRegex = regexp.MustCompile(`\$\([a-z\.\d]+\)`)
var partRegex = regexp.MustCompile(`[a-z\d]+`)

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
			C: cmdlist[0],
		}
		if len(cmdlist) > 1 {
			c.SubC = cmdlist[1:]
		}
		c.RawCmd = match
		// lazy fix to avoid out of range error
		c.SubC = append(c.SubC, "")
		c.SubC = append(c.SubC, "")
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
		line = strings.Replace(line, c.RawCmd, c.Outcome, 1)
	}
	return line
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
