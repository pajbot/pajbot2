package utils

import (
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
)

func TestParseTwitchDuration(t *testing.T) {
	c := qt.New(t)

	type testRun struct {
		input            string
		inputDefaultUnit time.Duration
		inputDefaultTime time.Duration
		expectedOutput   time.Duration
	}

	tests := []testRun{
		{
			input:            "1m",
			inputDefaultUnit: time.Second,
			inputDefaultTime: 10 * time.Minute,
			expectedOutput:   1 * time.Minute,
		},
		{
			input:            "1",
			inputDefaultUnit: time.Second,
			inputDefaultTime: 10 * time.Minute,
			expectedOutput:   1 * time.Second,
		},
		{
			input:            "1h",
			inputDefaultUnit: time.Second,
			inputDefaultTime: 10 * time.Minute,
			expectedOutput:   1 * time.Hour,
		},
		{
			input:            "1d",
			inputDefaultUnit: time.Second,
			inputDefaultTime: 10 * time.Minute,
			expectedOutput:   24 * time.Hour,
		},
	}

	for _, test := range tests {
		output := ParseTwitchDuration(test.input, test.inputDefaultUnit, test.inputDefaultTime)
		c.Assert(test.expectedOutput, qt.Equals, output)
	}
}
