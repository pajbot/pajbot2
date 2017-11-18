package command

import (
	"log"
	"strings"
	"time"

	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/helper"
)

// Handler is a helper struct/class to handle commands
type Handler struct {
	commands []Command
}

// AddCommand adds the given command to the list of commands
func (h *Handler) AddCommand(c Command) {
	bc := c.GetBaseCommand()
	// default values for now
	bc.Cooldown = time.Second * 5
	bc.UserCooldown = time.Second * 15
	h.commands = append(h.commands, c)
}

// RemoveCommand removes the given command from the list if commands
func (h *Handler) RemoveCommand(bc *BaseCommand) {
	for i, fC := range h.commands {
		fBc := fC.GetBaseCommand()
		if fBc.ID == bc.ID {
			h.commands = append(h.commands[:i], h.commands[i+1:]...)
			break
		}
	}
}

// GetTriggeredCommand returns the command that was triggered
func (h *Handler) GetTriggeredCommand(text string) Command {
	m := helper.GetTriggers(text)
	trigger := m[0]

	for _, cmd := range h.commands {
		if triggered, c := cmd.IsTriggered(trigger, m, 0); triggered {
			return c
		}
	}
	return nil
}

// Check will check a message if it should trigger a command
func (h *Handler) Check(b *bot.Bot, msg *common.Msg, action *bot.Action) error {
	if len(msg.Text) == 0 {
		return nil
	}

	m := helper.GetTriggers(msg.Text)

	if msg.Text[0] != '!' {
		return nil
	}
	c := h.GetTriggeredCommand(msg.Text)
	if c != nil {
		// Is the user high level enough to use this command?
		bc := c.GetBaseCommand()
		if bc.Level > msg.User.Level {
			log.Printf("%s tried to use %s, which requires level %d (he is level %d)",
				msg.User.DisplayName, strings.Join(m, " "), bc.Level, msg.User.Level)
			return nil
		}
		if c.OnCooldown(&msg.User) {
			log.Println("COMMAND IS ON COOLDOWN XDDD")
			return nil
		}
		// TODO: Get response first, and skip if the response is nil or something of that sort
		r := c.Run(b, msg, action)
		if r != "" {
			args := strings.Split(msg.Text, " ")
			if len(args) > 1 {
				msg.Args = args[1:]
				log.Println(msg.Args)
			}
			b.SayFormat(r, msg)
		}
	}
	return nil
}
