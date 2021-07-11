package banned_names

import (
	"regexp"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/modules"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
)

func init() {
	modules.Register("banned_names", func() pkg.ModuleSpec {
		// TODO: Make configurable
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

		return modules.NewSpec("banned_names", "Banned names", false, func(b *mbase.Base) pkg.Module {
			return newBannedNames(b, badUsernames)
		})
	})
}

type bannedNames struct {
	mbase.Base

	badUsernames []*regexp.Regexp
}

func newBannedNames(b *mbase.Base, badUsernames []*regexp.Regexp) pkg.Module {
	return &bannedNames{
		Base: *b,

		badUsernames: badUsernames,
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
