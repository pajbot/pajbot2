package command

import "testing"

func TestIsTriggered(t *testing.T) {
	var tests = []struct {
		command  Command
		trigger  string
		expected bool
	}{
		{
			command: Command{
				Triggers: []string{
					"a",
					"b",
				},
			},
			trigger:  "xd",
			expected: false,
		},
		{
			command: Command{
				Triggers: []string{},
			},
			trigger:  "xd",
			expected: false,
		},
		{
			command: Command{
				Triggers: []string{},
			},
			trigger:  "",
			expected: false,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
				},
			},
			trigger:  "",
			expected: false,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
				},
			},
			trigger:  "testa",
			expected: false,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
				},
			},
			trigger:  "atest",
			expected: false,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
				},
			},
			trigger:  "!test", // the !-parsing is handled by the module
			expected: false,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
				},
			},
			trigger:  "test",
			expected: true,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
					"abc",
				},
			},
			trigger:  "test",
			expected: true,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
					"abc",
				},
			},
			trigger:  "abc",
			expected: true,
		},
		{
			command: Command{
				Triggers: []string{
					"test",
					"abc",
				},
			},
			trigger:  "abcd",
			expected: false,
		},
	}

	for _, tt := range tests {
		res := tt.command.IsTriggered(tt.trigger)

		if res != tt.expected {
			if tt.expected {
				t.Errorf("%s failed triggering %s", tt.trigger, tt.command.Triggers)
			} else {
				t.Errorf("%s triggered %s when it shouldn't have", tt.trigger, tt.command.Triggers)
			}
		}
	}
}
