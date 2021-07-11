package modules

import (
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/commands"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/utils"
)

func init() {
	Register("value", func() pkg.ModuleSpec {
		return &Spec{
			id:               "value",
			name:             "Value",
			enabledByDefault: false,

			maker: newValue,

			parameters: map[string]pkg.ModuleParameterSpec{
				"A": func() pkg.ModuleParameter {
					return NewFloatParameter(ParameterSpec{
						Description:  "A kjdfghk jdfhgkj dfg",
						DefaultValue: float32(3.0),
					})
				},
				"B": func() pkg.ModuleParameter {
					return NewFloatParameter(ParameterSpec{
						Description:  "Bdfkgjh dkfjgh sdfgkkkk",
						DefaultValue: float32(6.0),
					})
				},
			},
		}
	})
}

type value struct {
	mbase.Base

	commands pkg.CommandsManager

	A float32
	B float32
}

type valueCmd struct {
	m     *value
	key   string
	value *float32
}

func (c *valueCmd) set(parts []string, event pkg.MessageEvent) pkg.Actions {
	if err := c.m.SetParameter(c.key, parts[1]); err != nil {
		return twitchactions.Mention(event.User, err.Error())
	}

	err := c.m.Save()
	if err != nil {
		return nil
	}
	return twitchactions.Mention(event.User, c.key+" set to "+utils.Float32ToString(*c.value))
}

func (c *valueCmd) get(parts []string, event pkg.MessageEvent) pkg.Actions {
	return twitchactions.Mention(event.User, c.key+" is "+utils.Float32ToString(*c.value))
}

func (c *valueCmd) Trigger(parts []string, event pkg.MessageEvent) pkg.Actions {
	if len(parts) >= 2 {
		return c.set(parts, event)
	}

	return c.get(parts, event)
}

func newValue(b *mbase.Base) pkg.Module {
	m := &value{
		Base: *b,

		commands: commands.NewCommands(),
	}

	m.Parameters()["A"].Link(&m.A)
	m.Parameters()["B"].Link(&m.B)

	m.commands.Register([]string{"!value-a"}, &valueCmd{m, "A", &m.A})
	m.commands.Register([]string{"!value-b"}, &valueCmd{m, "B", &m.B})

	return m
}

func (m *value) OnMessage(event pkg.MessageEvent) pkg.Actions {
	return m.commands.OnMessage(event)
}
