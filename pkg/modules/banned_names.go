package modules

import (
	"regexp"

	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

func init() {
	Register("banned_names", func() pkg.ModuleSpec {
		badUsernames := []*regexp.Regexp{
			regexp.MustCompile(`tos_is_trash\d+`),
			regexp.MustCompile(`trash_is_the_tos\d+`),
			regexp.MustCompile(`terms_of_service_uncool\d+`),
			regexp.MustCompile(`tos_i_love_mods_no_toxic\d+`),
			regexp.MustCompile(`^kemper.+`),
			regexp.MustCompile(`^pudele\d+`),
			regexp.MustCompile(`^ninjal0ver\d+`),
			regexp.MustCompile(`^trihard_account_\d+`),
			regexp.MustCompile(`^h[il1]erot[il1]tan.+`),
		}

		return &moduleSpec{
			id:   "banned_names",
			name: "Banned names",
			maker: func(b mbase.Base) pkg.Module {
				return newBannedNames(b, badUsernames)
			},

			enabledByDefault: false,
		}
	})
}

type bannedNames struct {
	mbase.Base

	badUsernames []*regexp.Regexp
}

func newBannedNames(b mbase.Base, badUsernames []*regexp.Regexp) pkg.Module {
	return &bannedNames{
		Base: b,
	}
}

func (m bannedNames) OnMessage(event pkg.MessageEvent) pkg.Actions {
	user := event.User

	usernameBytes := []byte(user.GetName())
	for _, badUsername := range m.badUsernames {
		if badUsername.Match(usernameBytes) {
			actions := &twitchactions.Actions{}
			actions.Ban(user).SetReason("Ban evasion")
			return actions
		}
	}

	return nil
}
