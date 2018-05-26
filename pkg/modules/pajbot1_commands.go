package modules

import (
	"log"
	"strings"

	twitch "github.com/gempir/go-twitch-irc"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/commands"
)

type Pajbot1Commands struct {
	server *server

	commands []*commands.Pajbot1Command

	Sender pkg.Channel
}

func NewPajbot1Commands(sender pkg.Channel) *Pajbot1Commands {
	return &Pajbot1Commands{
		server: &_server,
		Sender: sender,
	}
}

func (m *Pajbot1Commands) loadPajbot1Commands() error {
	const queryF = `SELECT level, action, command, delay_all, delay_user, enabled, cost, can_execute_with_whisper, sub_only, mod_only, tokens_cost FROM tb_command`

	session := m.server.oldSession

	stmt, err := session.Prepare(queryF)
	if err != nil {
		return err
	}

	rows, err := stmt.Query()
	if err != nil {
		return err
	}

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

func (m Pajbot1Commands) OnMessage(channel string, user twitch.User, message twitch.Message) error {
	if channel != "snusbot" {
		return nil
	}

	parts := strings.Split(message.Text, " ")
	if len(parts) == 0 {
		return nil
	}

	for _, command := range m.commands {
		if command.IsTriggered(parts) {
			err := command.Trigger(channel, user, parts, m.Sender)
			if err != nil {
				return err
			}
			log.Println("Triggered command!")
			log.Println(command.Action)
		}
	}

	return nil
}
