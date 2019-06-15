package modules

import (
	"fmt"
	"strings"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/utils"
)

func init() {
	Register("value", func() pkg.ModuleSpec {
		return &moduleSpec{
			id:               "value",
			name:             "Value",
			enabledByDefault: false,

			maker: newValue,

			parameters: map[string]pkg.ModuleParameterSpec{
				"A": func() pkg.ModuleParameter {
					return newFloatParameter(parameterSpec{
						Description:  "A kjdfghk jdfhgkj dfg",
						DefaultValue: float32(3.0),
					})
				},
				"B": func() pkg.ModuleParameter {
					return newFloatParameter(parameterSpec{
						Description:  "Bdfkgjh dkfjgh sdfgkkkk",
						DefaultValue: float32(6.0),
					})
				},
			},
		}
	})
}

type value struct {
	base

	A float32
	B float32
}

func newValue(b base) pkg.Module {
	m := &value{
		base: b,
	}

	m.parameters["A"].Link(&m.A)
	m.parameters["B"].Link(&m.B)

	return m
}

func (m *value) OnMessage(event pkg.MessageEvent) pkg.Actions {
	fmt.Println("VALUE MODULE ON MESSAGE")
	user := event.User
	message := event.Message

	if strings.HasPrefix(message.GetText(), "!") {
		parts := strings.Split(message.GetText(), " ")
		if parts[0] == "!value-a" {
			if len(parts) >= 2 {
				if err := m.setParameter("A", parts[1]); err != nil {
					return twitchactions.Mention(user, err.Error())
				}

				err := saveModule(m)
				if err != nil {
					fmt.Println("ERROR SAVING:", err)
				}
				return twitchactions.Mention(user, "A set to "+utils.Float32ToString(m.A))
			}

			return twitchactions.Mention(user, "A is "+utils.Float32ToString(m.A))
		}

		if parts[0] == "!value-b" {
			if len(parts) >= 2 {
				if err := m.setParameter("B", parts[1]); err != nil {
					return twitchactions.Mention(user, err.Error())
				}

				err := saveModule(m)
				if err != nil {
					fmt.Println("ERROR SAVING:", err)
				}
				return twitchactions.Mention(user, "B set to "+utils.Float32ToString(m.B))
			}

			return twitchactions.Mention(user, "B is "+utils.Float32ToString(m.B))
		}
	}

	return nil
}
