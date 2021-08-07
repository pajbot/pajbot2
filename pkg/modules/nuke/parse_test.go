package nuke

import (
	"regexp"
	"testing"
	"time"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseNukeParameters(t *testing.T) {
	c := qt.New(t)

	parser := &NukeParameterParser{}

	type testRun struct {
		input  []string
		output NukeParameters
	}

	tests := []testRun{
		{
			input: []string{"!nuke", "test", "1m", "10m"},
			output: NukeParameters{
				Phrase:           "test",
				RegexPhrase:      nil,
				ScrollbackLength: time.Minute * 1,
				TimeoutDuration:  time.Minute * 10,
			},
		},
		{
			input: []string{"!nuke", "test", "foo", "bar", "1m", "10m"},
			output: NukeParameters{
				Phrase:           "test foo bar",
				RegexPhrase:      nil,
				ScrollbackLength: time.Minute * 1,
				TimeoutDuration:  time.Minute * 10,
			},
		},
		{
			input: []string{"!nuke", "/test", "foo", "bar/", "1m", "10m"},
			output: NukeParameters{
				Phrase:           "/test foo bar/",
				RegexPhrase:      regexp.MustCompile("/test foo bar/"),
				ScrollbackLength: time.Minute * 1,
				TimeoutDuration:  time.Minute * 10,
			},
		},
		{
			input: []string{"!nuke", "/test/", "1m", "10m"},
			output: NukeParameters{
				Phrase:           "/test/",
				RegexPhrase:      regexp.MustCompile("test"),
				ScrollbackLength: time.Minute * 1,
				TimeoutDuration:  time.Minute * 10,
			},
		},
		{
			input: []string{"!nuke", "/test/", "1m", "1h"},
			output: NukeParameters{
				Phrase:           "/test/",
				RegexPhrase:      regexp.MustCompile("test"),
				ScrollbackLength: time.Minute * 1,
				TimeoutDuration:  time.Hour * 1,
			},
		},
		{
			input: []string{"!nuke", "/test/", "1m", "7h"},
			output: NukeParameters{
				Phrase:           "/test/",
				RegexPhrase:      regexp.MustCompile("test"),
				ScrollbackLength: time.Minute * 1,
				TimeoutDuration:  time.Hour * 7,
			},
		},
	}

	for _, test := range tests {
		output, _ := parser.ParseNukeParameters(test.input)
		c.Assert(test.output, qt.CmpEquals(cmpopts.IgnoreFields(NukeParameters{}, "RegexPhrase")), *output)
		if test.output.RegexPhrase == nil {
			c.Assert(output.RegexPhrase, qt.IsNil)
		} else {
			c.Assert(output.RegexPhrase, qt.Not(qt.IsNil))
		}
	}
}
