package command

import (
	"testing"

	"github.com/pajlada/pajbot2/helper"
)

func TestIsTriggered(t *testing.T) {
	var tests = []struct {
		command  Command
		message  string
		expected bool
	}{
		{
			command: &TextCommand{
				Triggers: []string{
					"a",
					"b",
				},
			},
			message:  "!xd",
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{},
			},
			message:  "!xd",
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{},
			},
			message:  "",
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
				},
			},
			message:  "",
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
				},
			},
			message:  "!testa",
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
				},
			},
			message:  "!atest",
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
				},
			},
			message:  "!!test", // the !-parsing is handled by the module
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
				},
			},
			message:  "!test",
			expected: true,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
					"abc",
				},
			},
			message:  "!test",
			expected: true,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
					"abc",
				},
			},
			message:  "!abc",
			expected: true,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
					"abc",
				},
			},
			message:  "!abcd",
			expected: false,
		},
		{
			command: &TextCommand{
				Triggers: []string{
					"test",
					"abc",
				},
			},
			message:  "!abcd LALALA",
			expected: false,
		},
	}

	for _, tt := range tests {
		m := helper.GetTriggers(tt.message)
		trigger := m[0]

		triggered, _ := tt.command.IsTriggered(trigger, m, 0)

		if triggered != tt.expected {
			if tt.expected {
				t.Errorf("%s failed triggering", tt.message)
			} else {
				t.Errorf("%s triggered when it shouldn't have", tt.message)
			}
		}
	}
}
