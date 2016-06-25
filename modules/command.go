package modules

import (
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
	commands []command.Command
}

/*
some ideas:
not sure how this would work but id like to keep things simple

$(user1) = usersource displayname
$(sender) = sender of msg
$(1) = raw arg from msg
$(bot) = bot
$(channel) = channel

users have methods, for example $(sender.points):
	.low = lowercase
	.points
	.level
	etc ...

bot has:
	bot.uptime
	bot.version
	etc...

channel has:
	channel.uptime
	channel.title
	channel.game
	channel.subs
	channel.name
	etc...

!points would look like this :
	"$(user1) has $(user1.points) points."

!uptime:
	"$(sender), $(channel.name) has been online for $(channel.uptime) PogChamp"

*/

// Ensure the module implements the interface properly
var _ Module = (*Command)(nil)

// Init initializes something
func (module *Command) Init(sql *sqlmanager.SQLManager) {
	// Fetch rows from pb_command
	rows, err := sql.Session.Query("SELECT * FROM pb_command")

	if err != nil {
		log.Fatal("Error fetching commands:", err)
	}

	for rows.Next() {
		var triggers string
		if err := rows.Scan(&triggers); err != nil {
			log.Fatal(err)
		}
		log.Debug("triggers: ", triggers)
		c := command.TextCommand{
			Triggers: strings.Split(triggers, "|"),
			Response: "pajaSWA xD",
		}
		module.commands = append(module.commands, &c)
	}

	xdCommand := command.TextCommand{
		Triggers: []string{
			"xd",
			"xdlol",
		},
		Response: "pajaSWA",
	}
	module.commands = append(module.commands, &xdCommand)
	testCommand := command.NestedCommand{
		Triggers: []string{
			"lul",
			"xdlul",
		},
		Commands: []command.Command{
			&xdCommand,
			&command.TextCommand{
				Triggers: []string{
					"a",
				},
				Response: "pajaSWA a ;P",
			},
			&command.TextCommand{
				Triggers: []string{
					"b",
				},
				Response: "pajaSWA b ;P",
			},
		},
	}
	module.commands = append(module.commands, &testCommand)
}

// Check xD
func (module *Command) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if len(msg.Text) == 0 {
		// Do nothing with empty messages
		return nil
	}

	m := helper.GetTriggers(msg.Text)
	trigger := m[0]

	if msg.Text[0] != '!' {
		return nil
	}
	for _, command := range module.commands {
		if triggered, c := command.IsTriggered(trigger, m, 0); triggered {
			// TODO: Get response first, and skip if the response is nil or something of that sort
			r := c.Run()
			if r != "" {
				action.Response = r
				action.Stop = true
			}
			return nil
		}
	}
	return nil
}
