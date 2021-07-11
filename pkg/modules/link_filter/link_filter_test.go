package link_filter

import (
	"testing"

	"github.com/pajbot/pajbot2/pkg/modules"
)

func TestLinkFilterUnmatches(t *testing.T) {
	spec, ok := modules.GetModuleSpec("link_filter")
	if !ok {
		t.Fatal("AAAAAA")
	}

	m := spec.Create(nil).(*LinkFilter)

	// these strings should not match as links
	tests := []string{
		"test.Bb5",
	}

	for _, link := range tests {
		triggered := m.checkMessage(link)
		if triggered {
			t.Errorf("%s is seen as a link while it should not be", link)
		}
	}
}

func TestLinkFilterMatches(t *testing.T) {
	spec, ok := GetModuleSpec("link_filter")
	if !ok {
		t.Fatal("AAAAAA")
	}

	m := spec.Create(nil).(*LinkFilter)

	// these strings should match as links
	tests := []string{
		"google.com",
		"mylovely.horse",
		"clips.twitch.tv",
		"google.nl",
		"google.co.uk",
		"google.fi",
		"twitch.tv",
		"pajlada.se",
		"test.bb",
		"https://google.com",
		"https://twitter.com",
	}

	for _, link := range tests {
		triggered := m.checkMessage(link)
		if !triggered {
			t.Errorf("%s is not seen as a link while it should be", link)
		}
	}
}
