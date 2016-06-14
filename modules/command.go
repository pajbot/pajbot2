package modules

import (
	"log"
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
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
		log.Println("triggers: ", triggers)
		c := command.Command{
			Triggers: strings.Split(triggers, "|"),
			Response: "pajaSWA xD",
		}
		module.commands = append(module.commands, c)
	}

	xdCommand := command.Command{
		Triggers: []string{
			"xd",
			"xdlol",
		},
		Response: "pajaSWA",
	}
	module.commands = append(module.commands, xdCommand)
}

// Check xD
func (module *Command) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	m := strings.Split(msg.Message, " ")
	trigger := strings.ToLower(m[0])
	if trigger[0] != '!' {
		return nil
	}
	trigger = trigger[1:]
	for _, command := range module.commands {
		if command.IsTriggered(trigger) {
			// TODO: Get response first, and skip if the response is nil or something of that sort
			action.Response = command.GetResponse()
			action.Stop = true
			return nil
		}
	}
	return nil
}
