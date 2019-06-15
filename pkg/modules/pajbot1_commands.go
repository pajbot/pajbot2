package modules

import (
	"log"
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
)

func init() {
	Register("pajbot1_commands", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:    "pajbot1_commands",
			name:  "pajbot1 commands",
			maker: newPajbot1Commands,
		}
	})
}

type Pajbot1Commands struct {
	base

	commands []*commands.Pajbot1Command
}

func newPajbot1Commands(b base) pkg.Module {
	m := &Pajbot1Commands{
		base: b,
	}

	m.loadPajbot1Commands()

	return m
}

func (m *Pajbot1Commands) loadPajbot1Commands() error {
	const queryF = `SELECT level, action, command, delay_all, delay_user, enabled, cost, can_execute_with_whisper, sub_only, mod_only, tokens_cost FROM tb_command`

	session := m.server.oldSession

	rows, err := session.Query(queryF)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var command commands.Pajbot1Command
		err = command.LoadScan(rows)
		if err != nil {
			return err
		}

		if !command.Enabled {
			continue
		}

		if command.PointCost > 0 || command.TokenCost > 0 {
			continue
		}

		m.commands = append(m.commands, &command)
	}

	return nil
}

func (m Pajbot1Commands) OnMessage(event pkg.MessageEvent) pkg.Actions {
	user := event.User
	message := event.Message

	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	for _, command := range m.commands {
		if command.IsTriggered(parts) {
			err := command.Trigger(m.bot, user, parts)
			if err != nil {
				return nil
			}
			log.Println("Triggered command!")
			log.Println(command.Action)
		}
	}

	return nil
}
