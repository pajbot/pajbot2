package commands_test

import (
	"testing"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	"github.com/pajbot/pajbot2/pkg/users"
)

type testSimpleCommand struct {
	n int
}

func (c *testSimpleCommand) Trigger(args []string, e pkg.MessageEvent) pkg.Actions {
	c.n += 1
	return nil
}

type testCustomCommand struct {
	testSimpleCommand

	cds map[string]bool
}

func (c *testCustomCommand) HasCooldown(user pkg.User) bool {
	if c.cds == nil {
		c.cds = make(map[string]bool)
	}

	return c.cds[user.GetName()]
}

func (c *testCustomCommand) AddCooldown(user pkg.User) {
	if c.cds == nil {
		c.cds = make(map[string]bool)
	}

	c.cds[user.GetName()] = true
}

type twitchMessage struct {
	msg string
}

func (m *twitchMessage) GetText() string {
	return m.msg
}
func (m *twitchMessage) SetText(text string) {
	m.msg = text
}

func (m *twitchMessage) GetTwitchReader() pkg.EmoteReader {
	return nil
}

func (m *twitchMessage) GetBTTVReader() pkg.EmoteReader {
	return nil
}

func (m *twitchMessage) AddBTTVEmote(emote pkg.Emote) {
}

func TestSimpleCommand(t *testing.T) {
	cmds := commands.NewCommands()

	cmd := &testSimpleCommand{}

	if cmd.n != 0 {
		t.Fatal("Command was not initialized with value 0")
	}

	cmds.Register2(5, []string{"!foo"}, cmd)

	msg := twitchMessage{
		msg: "!foo bar",
	}

	event := pkg.MessageEvent{
		BaseEvent: pkg.BaseEvent{
			UserStore: nil,
		},

		User:    nil,
		Message: &msg,
		Channel: nil,
	}

	cmds.OnMessage(event)

	if cmd.n != 1 {
		t.Fatal("Command should have been executed once")
	}
}

func TestCustomCommand(t *testing.T) {
	cmds := commands.NewCommands()

	cmd := &testCustomCommand{}

	if cmd.n != 0 {
		t.Fatal("Command was not initialized with value 0")
	}

	cmds.Register2(5, []string{"!foo"}, cmd)

	msg := twitchMessage{
		msg: "!foo bar",
	}

	user := users.NewSimpleTwitchUser("11148817", "pajlada")

	if cmd.HasCooldown(user) {
		t.Fatal("User should not have a cooldown if the command has not been run")
	}

	event := pkg.MessageEvent{
		BaseEvent: pkg.BaseEvent{
			UserStore: nil,
		},

		User:    user,
		Message: &msg,
		Channel: nil,
	}

	cmds.OnMessage(event)

	if cmd.n != 1 {
		t.Fatal("Command should have been executed once")
	}

	if !cmd.HasCooldown(user) {
		t.Fatal("User should have cooldown now that the command has been run")
	}
}
