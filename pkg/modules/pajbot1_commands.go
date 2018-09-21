package modules

import (
	"log"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/commands"
)

type Pajbot1Commands struct {
	server *server

	commands []*commands.Pajbot1Command

	Sender pkg.Sender
}

func NewPajbot1Commands(sender pkg.Sender) *Pajbot1Commands {
	return &Pajbot1Commands{
		server: &_server,
		Sender: sender,
	}
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

func (m *Pajbot1Commands) Register() error {

	err := m.loadPajbot1Commands()
	if err != nil {
		return err
	}

	return nil
}

func (m Pajbot1Commands) Name() string {
	return "Pajbot1Commands"
}

func (m Pajbot1Commands) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m Pajbot1Commands) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if source.GetChannel() != "snusbot" {
		return nil
	}

	parts := strings.Split(message.GetText(), " ")
	if len(parts) == 0 {
		return nil
	}

	for _, command := range m.commands {
		if command.IsTriggered(parts) {
			err := command.Trigger(source, user, parts, m.Sender)
			if err != nil {
				return err
			}
			log.Println("Triggered command!")
			log.Println(command.Action)
		}
	}

	return nil
}
