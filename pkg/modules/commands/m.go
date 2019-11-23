// Package mcommands describes the "Commands" module
//
// Add command
// Chat: !pb2addcmd trigger response
// Web: not implemented
// Code: not implemented
//
// Add trigger to command
// Chat: !pb2addtrigger command_id new_trigger
// Web: not implemented
//
// Remove trigger from command
// Chat: !pb2removetrigger command_id trigger
// Web: not implemented
//
// Edit command
// Chat: not implemented
// Web: not implemented
//
// Edit trigger
// Chat: not implemented
// Web: not implemented
//
// Remove trigger
// Chat: not implemented
// Web: not implemented
//
// Remove command
// Chat: not implemented
// Web: not implemented
package mcommands

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/pajbot/pajbot2/internal/utils"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
)

func init() {
	modules.Register("commands", func() pkg.ModuleSpec {
		return modules.NewSpec(
			"commands",
			"Commands",
			false,
			newCommands,
		)
	})
}

type CommandsModule struct {
	mbase.Base

	commands pkg.CommandsManager
}

func newCommands(b mbase.Base) pkg.Module {
	m := &CommandsModule{
		Base: b,

		commands: commands.NewCommands(),
	}

	m.commands.Register([]string{"!pb2addcmd"}, newAddCmd(m))
	m.commands.Register([]string{"!pb2addtrigger"}, newCmdAddTrigger(m))
	m.commands.Register([]string{"!pb2removetrigger"}, newCmdRemoveTrigger(m))

	err := m.loadCommands()
	if err != nil {
		fmt.Println("Error loading ocmmands:", err)
	}

	return m
}

func (m *CommandsModule) loadCommands() error {
	const querySelect = `
SELECT
	id, response, command_trigger.trigger
FROM
	command
LEFT JOIN
	command_trigger
ON
	command.id=command_trigger.command_id
WHERE
	command.bot_channel_id=$1
	`
	rows, err := m.SQL.Query(querySelect, m.BotChannel().DatabaseID())
	if err != nil {
		return err
	}
	cmds := map[int64]*cmdData{}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var response, trigger string
		err := rows.Scan(&id, &response, &trigger)
		if err != nil {
			return err
		}
		cmd := cmds[id]
		if cmd == nil {
			cmd = &cmdData{}
		}
		cmd.id = id
		cmd.response = response
		cmd.triggers = append(cmd.triggers, trigger)
		cmds[id] = cmd
	}

	for _, cmd := range cmds {
		fmt.Println("Registering triggers", cmd.triggers, "for response", cmd.response)
		m.commands.Register2(cmd.id, cmd.triggers, newTextResponseCommand(cmd.response))
	}

	fmt.Println(cmds)

	return nil
}

func (m *CommandsModule) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

func (m *CommandsModule) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}

func (m *CommandsModule) addToDB(trigger, response string) error {
	err := utils.WithTransaction(m.SQL, func(tx *sql.Tx) error {
		const queryInsertCommand = `
		INSERT INTO
			command
			(bot_channel_id, action_type, response)
		VALUES
			($1, 'text_response', $2)
		RETURNING id
			`

		const queryTrigger = `
		INSERT INTO
			command_trigger
			(command_id, trigger)
		VALUES
			($1, $2)
			`

		var id int64

		row := tx.QueryRow(queryInsertCommand, m.BotChannel().DatabaseID(), response)
		err := row.Scan(&id)
		if err != nil {
			return err
		}

		_, err = tx.Exec(queryTrigger, id, trigger)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	m.commands.Register([]string{trigger}, newTextResponseCommand(response))

	return nil
}

func (m *CommandsModule) addTrigger(commandID int64, newTrigger string) error {
	err := utils.WithTransaction(m.SQL, func(tx *sql.Tx) error {
		const queryTrigger = `
		INSERT INTO
			command_trigger
			(command_id, trigger)
		VALUES
			($1, $2)
			`

		_, err := tx.Exec(queryTrigger, commandID, newTrigger)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	cmd := m.commands.FindByCommandID(commandID)
	if cmd == nil {
		return errors.New("command with id does not exist internally??")
	}

	m.commands.Register([]string{newTrigger}, cmd)

	return nil
}

func (m *CommandsModule) removeTrigger(commandID int64, trigger string) error {
	const queryRemoveTrigger = `
	DELETE FROM
		command_trigger
	WHERE
		command_id=$1 AND
		trigger=$2
		`

	res, err := m.SQL.Exec(queryRemoveTrigger, commandID, trigger)
	if err != nil {
		return err
	}

	if ra, err := res.RowsAffected(); err != nil {
		return err
	} else if ra == 0 {
		return errors.New("no sql trigger to delete")
	}

	m.commands.DeregisterAliases([]string{trigger})

	return nil
}

type cmdData struct {
	id       int64
	response string
	triggers []string
}
