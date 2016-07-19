package modules

import (
	"strings"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/command"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/helper"
)

/*
Admin xD
*/
type Admin struct {
	commandHandler command.Handler
}

// Ensure the module implements the interface properly
var _ Module = (*Admin)(nil)

func cmdJoinChannel(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	m := helper.GetTriggersN(msg.Text, 2)

	if len(m) < 1 {
		b.Say("Usage: !admin join forsenlol")
		// Not enough arguments
		return
	}

	newChannel := strings.ToLower(strings.Replace(m[0], "#", "", -1))

	// Fetch all existing channels and see if this channel is already there
	channels, err := common.FetchAllChannels(b.SQL)
	if err != nil {
		b.Say("Errors fetching existing channels")
		return
	}

	isChannelNew := true
	for _, channel := range channels {
		if channel.Name == newChannel {
			// We already know of this channel
			isChannelNew = false

			// Is the channel enabled?
			if channel.Enabled {
				// We are already in this channel
				b.Sayf("We are already in the channel %s", newChannel)
				return
			}

			// We've been in this channel before, but not currently
			channel.Enabled = true
			b.Sayf("Rejoining %s", newChannel)
			channel.SQLSetEnabled(b.SQL, 1)

			break
		}
	}

	if isChannelNew {
		// Insert new channel
		c := &common.Channel{
			Name: newChannel,
		}
		c.InsertNewToSQL(b.SQL)
	}

	b.Join <- newChannel
}

func cmdLeaveChannel(b *bot.Bot, msg *common.Msg, action *bot.Action) {
	m := helper.GetTriggersN(msg.Text, 2)

	if len(m) < 1 {
		b.Say("Usage: !admin leave forsenlol")
		// Not enough arguments
		return
	}

	newChannelName := strings.ToLower(strings.Replace(m[0], "#", "", -1))

	// Fetch all existing channels and see if this channel is already there
	channels, err := common.FetchAllChannels(b.SQL)
	if err != nil {
		b.Say("Errors fetching existing channels")
		return
	}

	for _, channel := range channels {
		if channel.Name == newChannelName {
			// Is the channel enabled?
			if channel.Enabled {
				// We've been in this channel before, but not currently
				channel.Enabled = false
				b.Sayf("Leaving %s", newChannelName)
				channel.SQLSetEnabled(b.SQL, 0)

				b.Join <- "PART " + newChannelName
				return
			}

			// We are not in this channel
			b.Sayf("We are not in the channel %s", newChannelName)
			return
		}
	}
}

// Init xD
func (module *Admin) Init(bot *bot.Bot) {
	testCommand := command.NestedCommand{
		BaseCommand: command.BaseCommand{
			Triggers: []string{
				"admin",
			},
			Level: 500,
		},
		Commands: []command.Command{
			&command.FuncCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"join",
						"joinchannel",
					},
					Level: 500,
				},
				Function: cmdJoinChannel,
			},
			&command.FuncCommand{
				BaseCommand: command.BaseCommand{
					Triggers: []string{
						"leave",
						"part",
						"leavechannel",
						"partchannel",
					},
					Level: 500,
				},
				Function: cmdLeaveChannel,
			},
		},
	}
	module.commandHandler.AddCommand(&testCommand)
}

// Check xD
func (module *Admin) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	return module.commandHandler.Check(b, msg, action)
}
