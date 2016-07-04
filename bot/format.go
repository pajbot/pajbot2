package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pajlada/pajbot2/common"
)

// Format xD
type Format struct {
	re    *regexp.Regexp
	cmdRe *regexp.Regexp
}

// command is not a good name, but idk what else to call it
type command struct {
	c       string
	subC    []string
	rawCmd  string
	outcome string
}

// InitFormatter compiles format regexes
func (bot *Bot) InitFormatter() *Format {
	return &Format{
		re:    regexp.MustCompile(`\$\([a-z\.]+\)`),
		cmdRe: regexp.MustCompile(`[a-z]+`),
	}
}

// Format formats the given line xD
func (bot *Bot) Format(line string, msg *common.Msg) string {
	// catch all errors until we have proper error handling
	defer func() {
		if r := recover(); r != nil {
			log.Error(r)
		}
	}()
	f := bot.Fmt
	fmtline, rawCommands := f.parseLine(line)
	for i := range rawCommands {
		bot.execCommand(&rawCommands[i], msg)
	}
	return f.formatCommands(fmtline, rawCommands)
}

func (bot *Bot) execCommand(cmd *command, msg *common.Msg) {
	switch cmd.c {
	case "source", "sender":
		cmd.outcome = bot.formatUser(&msg.User, cmd.subC)
	case "user":
		if msg.Args != nil {
			if bot.Redis.IsValidUser(msg.Channel, msg.Args[0]) {
				user := bot.Redis.LoadUser(msg.Channel, msg.Args[0])
				cmd.outcome = bot.formatUser(&user, cmd.subC)
				return
			}
		}
		cmd.outcome = bot.formatUser(&msg.User, cmd.subC)
	}
}

func (bot *Bot) formatUser(user *common.User, cmds []string) string {
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
	default:
		return user.DisplayName
	}
}

func (f *Format) formatCommands(line string, cmds []command) string {
	log.Debug(line)
	log.Debug(cmds)
	for _, c := range cmds {
		line = strings.Replace(line, c.rawCmd, c.outcome, 1)
	}
	log.Debug(line)
	return line
}

func (f *Format) parseLine(line string) (string, []command) {
	log.Debug(line)
	matches := f.re.FindAllString(line, -1)
	log.Debug(matches)
	var cmds []command
	for _, match := range matches {
		cmdlist := f.cmdRe.FindAllString(match, -1)
		c := command{
			c: cmdlist[0],
		}
		if len(cmdlist) > 1 {
			c.subC = cmdlist[1:]
		}
		c.rawCmd = match
		cmds = append(cmds, c)
	}
	log.Debug(cmds, line)
	return line, cmds
}
