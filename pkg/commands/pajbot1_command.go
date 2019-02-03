package commands

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

const commandPrefix = "!"

// Pajbot1Command is a command loaded from the old pajbot1 database
type Pajbot1Command struct {
	Level     int
	Action    string
	Triggers  []string
	Enabled   bool
	PointCost int

	GlobalCooldown int
	UserCooldown   int

	CanExecuteWithWhisper bool
	SubOnly               bool
	ModOnly               bool

	TokenCost int

	// If a command has possible userdata, it means we should run it through our banphrase filter before printing it
	HasUserdata bool
}

type pajbot1CommandAction struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (c *Pajbot1Command) LoadScan(rows *sql.Rows) error {
	var actionString []byte
	var commandString string
	const queryF = `SELECT level, action, command, delay_all, delay_user, enabled, cost, can_execute_with_whisper, sub_only, mod_only, tokens_cost FROM tb_command`
	err := rows.Scan(&c.Level, &actionString, &commandString, &c.GlobalCooldown, &c.UserCooldown, &c.Enabled, &c.PointCost, &c.CanExecuteWithWhisper, &c.SubOnly, &c.ModOnly, &c.TokenCost)
	if err != nil {
		return err
	}

	var action pajbot1CommandAction

	err = json.Unmarshal(actionString, &action)
	if err != nil {
		return err
	}

	c.Triggers = strings.Split(commandString, "|")

	for key, t := range c.Triggers {
		c.Triggers[key] = commandPrefix + t
	}

	c.Action = action.Message

	return nil
}

func (c *Pajbot1Command) IsTriggered(parts []string) bool {
	for _, trigger := range c.Triggers {
		if strings.ToLower(parts[0]) == trigger {
			return true
		}
	}

	return false
}

func (c *Pajbot1Command) Trigger(bot pkg.BotChannel, user pkg.User, parts []string) error {
	bot.Say(c.Action)

	return nil
}
