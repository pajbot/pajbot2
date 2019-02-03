package modules

import (
	"regexp"

	"github.com/pajlada/pajbot2/pkg"
)

type bannedNames struct {
	botChannel pkg.BotChannel

	server *server

	badUsernames []*regexp.Regexp
}

func newBannedNames() pkg.Module {
	return &bannedNames{
		server: &_server,
	}
}

var bannedNamesSpec = moduleSpec{
	id:    "banned_names",
	name:  "Banned names",
	maker: newBannedNames,

	enabledByDefault: false,
}

func (m *bannedNames) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`tos_is_trash\d+`))
	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`trash_is_the_tos\d+`))
	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`terms_of_service_uncool\d+`))
	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`tos_i_love_mods_no_toxic\d+`))
	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`^kemper.+`))
	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`^pudele\d+`))
	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`^ninjal0ver\d+`))
	m.badUsernames = append(m.badUsernames, regexp.MustCompile(`^trihard_account_\d+`))

	return nil
}

func (m *bannedNames) Disable() error {
	return nil
}

func (m *bannedNames) Spec() pkg.ModuleSpec {
	return &bannedNamesSpec
}

func (m *bannedNames) BotChannel() pkg.BotChannel {
	return m.botChannel
}

func (m bannedNames) OnWhisper(bot pkg.BotChannel, source pkg.User, message pkg.Message) error {
	return nil
}

func (m bannedNames) OnMessage(bot pkg.BotChannel, user pkg.User, message pkg.Message, action pkg.Action) error {
	usernameBytes := []byte(user.GetName())
	for _, badUsername := range m.badUsernames {
		if badUsername.Match(usernameBytes) {
			action.Set(pkg.Ban{"Ban evasion"})
			return nil
		}
	}

	return nil
}
