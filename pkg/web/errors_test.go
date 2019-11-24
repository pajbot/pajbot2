package web

import "testing"

func TestIsValidUserName(t *testing.T) {
	good := []string{
		"asd",
		"pajlada",
		"randers",
		"karl_kons",
		"testaccount_420",
		"bajlada",
	}
	bad := []string{
		"lol hi",
		"!@$%!@$!@$!@$",
		" pajlada",
	}
	for _, name := range good {
		if !isValidUserName(name) {
			t.Fatalf("%s is not considered valid, while it should be", name)
		}
	}

	for _, name := range bad {
		if isValidUserName(name) {
			t.Fatalf("%s is considered valid, while it shouldn't be", name)
		}
	}
}
