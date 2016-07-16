package modules

import (
	"database/sql"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/helper"
	"github.com/pajlada/pajbot2/sqlmanager"
)

/*
Command xD
*/
type Command struct {
	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*Command)(nil)

func (module *Command) loadCommands(sql *sqlmanager.SQLManager, channel common.Channel) int {
	// Fetch rows from pb_command
	rows, err := sql.Session.Query("SELECT id, channel_id, triggers, response, response_type FROM pb_command")

	if err != nil {
		log.Error("Error fetching commands:", err)
		return 0
	}

	return module.readCommands(rows)
}

// loadCommand loads a command with a given ID
func (module *Command) loadCommand(sql *sqlmanager.SQLManager, commandID int64) int {
	// Fetch rows from pb_command
	rows, err := sql.Session.Query("SELECT id, channel_id, triggers, response, response_type FROM pb_command WHERE `id`=?", commandID)

	if err != nil {
		log.Error("Error fetching commands:", err)
		return 0
	}

	return module.readCommands(rows)
}

func (module *Command) readCommands(rows *sql.Rows) int {
	numCommands := 0

	for rows.Next() {
		c := command.ReadSQLCommand(rows)
		if c != nil {
			module.commandHandler.AddCommand(c)
			numCommands++
		}
	}
	return numCommands
}

// Init initializes something
func (module *Command) Init(bot *bot.Bot) {
	module.loadCommands(bot.SQL, bot.Channel)

	xdCommand := command.TextCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"xd",
				"xdlol",
			},
		},
		Response: "pajaSWA",
	}
	module.commandHandler.AddCommand(&xdCommand)
	testCommand := command.NestedCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"lul",
				"xdlul",
			},
		},
		Commands: []command.Command{
			&xdCommand,
			&command.TextCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"a",
					},
				},
				Response: "pajaSWA a ;P",
			},
			&command.TextCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"b",
					},
				},
				Response: "pajaSWA b ;P",
			},
		},
	}
	module.commandHandler.AddCommand(&testCommand)

	// Temporary !admin prefix while it's in development
	adminCommand := &command.NestedCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"admin",
			},
		},
		Commands: []command.Command{
			&command.NestedCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"add",
					},
				},
				Commands: []command.Command{
					&command.FuncCommand{
						BaseCommand: command.BaseCommand{
							Triggers: []string{
								"command",
							},
							Level: 500,
						},
						Function: module.createCommand,
					},
				},
			},
			&command.NestedCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"remove",
					},
				},
				Commands: []command.Command{
					&command.FuncCommand{
						BaseCommand: command.BaseCommand{
							Triggers: []string{
								"command",
							},
							Level: 500,
						},
						Function: module.removeCommand,
					},
				},
			},
			&xdCommand,
			&command.TextCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"a",
					},
				},
				Response: "pajaSWA a ;P",
			},
			&command.TextCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"b",
					},
				},
				Response: "pajaSWA b ;P",
			},
		},
	}
	module.commandHandler.AddCommand(adminCommand)
}

func (module *Command) createCommand(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	// Change to 2 when we remove the !admin prefix
	const triggerLength = 3
	const usageFormat = "Usage: !%s !trigger response"
	triggers := helper.GetTriggersKC(msg.Text)
	if len(triggers) < triggerLength {
		b.Say("Missing arguments")
		return
	}

	if len(triggers) < triggerLength+2 {
		b.Sayf(usageFormat, strings.Join(triggers[:triggerLength], " "))
		return
	}
	// TODO: use an argument parser so we can have --arguments like --silent and --reply --me --cd=0
	arguments := triggers[triggerLength:]

	triggerString := strings.Replace(strings.ToLower(arguments[0]), "!", "", -1)
	if len(triggerString) == 0 {
		b.Sayf(usageFormat, strings.Join(triggers[:triggerLength], " "))
		return
	}
	var triggerList []string
	for _, t := range strings.Split(triggerString, "|") {
		add := true
		if len(t) > 0 {
			for _, eT := range triggerList {
				if t == eT {
					add = false
					break
				}
			}
			if add {
				triggerList = append(triggerList, t)
			}
		}
	}
	triggerString = strings.Join(triggerList, "|")

	response := arguments[1:]

	// See if any of the aliases we want to use is already in use
	for _, trigger := range triggerList {
		c := module.commandHandler.GetTriggeredCommand("!" + trigger)
		if c != nil {
			b.Sayf("Command !%s is already in use.", trigger)
			return
		}
	}

	sqlCommand := command.SQLCommand{
		ChannelID: 1, // XXX
		Triggers:  triggerString,
		Response:  strings.Join(response, " "),
	}

	commandID := sqlCommand.Insert(b.SQL.Session)
	addedCommands := module.loadCommand(b.SQL, commandID)
	if addedCommands == 1 {
		b.Sayf("Successfully added command with triggers %s", triggerString)
	} else {
		b.Sayf("Something went wrong when adding the command, %d commands were added ???", addedCommands)
	}
}

func (module *Command) removeCommand(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	// Change to 2 when we remove the !admin prefix
	const triggerLength = 3
	const usageFormat = "Usage: !%s !command"
	triggers := helper.GetTriggersKC(msg.Text)
	if len(triggers) < triggerLength {
		b.Say("Missing arguments")
		return
	}

	if len(triggers) < triggerLength+1 {
		b.Sayf(usageFormat, strings.Join(triggers[:triggerLength], " "))
		return
	}
	// TODO: use an argument parser so we can have --arguments like --silent and --reply --me --cd=0
	arguments := triggers[triggerLength:]

	trigger := strings.Replace(strings.ToLower(arguments[0]), "!", "", -1)

	c := module.commandHandler.GetTriggeredCommand("!" + trigger)
	if c == nil {
		b.Sayf("No command with trigger !%s", trigger)
		return
	}

	bc := c.GetBaseCommand()

	sqlCommand := command.SQLCommand{
		ID: bc.ID,
	}

	// Delete from DB
	err := sqlCommand.Delete(b.SQL.Session)
	if err != nil {
		b.Sayf("Error deleting command: %s", err)
	} else {
		b.Sayf("Successfully deleted command with trigger !%s", trigger)
	}

	// Delete from slice
	module.commandHandler.RemoveCommand(bc)
}

// Check xD
func (module *Command) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}
