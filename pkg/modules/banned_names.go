package modules

import (
	"regexp"

	"github.com/pajlada/pajbot2/pkg"
)

type BannedNames struct {
	server *server

	badUsernames []*regexp.Regexp
}

func NewBannedNames() *BannedNames {
	return &BannedNames{
		server: &_server,
	}
}

func (m *BannedNames) Register() error {

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

func (m BannedNames) Name() string {
	return "BannedNames"
}

func (m BannedNames) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m BannedNames) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if source.GetChannel() != "forsen" {
		return nil
	}

	usernameBytes := []byte(user.GetName())
	for _, badUsername := range m.badUsernames {
		if badUsername.Match(usernameBytes) {
			action.Set(pkg.Ban{"Ban evasion"})
			return nil
		}
	}

	return nil
}
