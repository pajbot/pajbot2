package modules

import (
	"fmt"
	"time"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/filters"
	"github.com/pajlada/pajbot2/pkg/utils"
)

type pajbot1BanphraseFilter struct {
	botChannel pkg.BotChannel

	server *server

	banphrases []pkg.Banphrase
}

func newPajbot1BanphraseFilter() pkg.Module {
	return &pajbot1BanphraseFilter{
		server: &_server,
	}
}

var pajbot1BanphraseSpec = moduleSpec{
	id:    "pajbot1_banphrase",
	name:  "pajbot1 banphrase",
	maker: newPajbot1BanphraseFilter,

	enabledByDefault: true,
}

func (m *pajbot1BanphraseFilter) addCustomBanphrase(phrase string) {
	m.banphrases = append(m.banphrases, &filters.Pajbot1Banphrase{
		ID:            -1,
		Name:          "Custom",
		Phrase:        phrase,
		Length:        600,
		Operator:      filters.OperatorContains,
		Permanent:     false,
		Warning:       false,
		Notify:        false,
		CaseSensitive: false,
		Enabled:       true,
		SubImmunity:   false,
		RemoveAccents: true,
	})
}

func (m *pajbot1BanphraseFilter) loadPajbot1Banphrases() error {
	const queryF = `SELECT * FROM tb_banphrase`

	session := m.server.oldSession

	rows, err := session.Query(queryF)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var bp filters.Pajbot1Banphrase
		err = bp.LoadScan(rows)
		if err != nil {
			return err
		}

		if bp.Enabled {
			// m.banphrases = append(m.banphrases, &bp)
		}
	}

	return nil
}

func (m *pajbot1BanphraseFilter) Initialize(botChannel pkg.BotChannel, settings []byte) error {
	m.botChannel = botChannel

	// hard-coded banphrases
	m.addCustomBanphrase("n!66ger")

	m.addCustomBanphrase("negro")
	m.addCustomBanphrase("negr0")
	m.addCustomBanphrase("n3gro")
	m.addCustomBanphrase("n3gr0")

	m.addCustomBanphrase("nij1gger")

	m.addCustomBanphrase("nijgger")
	m.addCustomBanphrase("nijggger")
	m.addCustomBanphrase("nijgggger")
	m.addCustomBanphrase("nijggggger")
	m.addCustomBanphrase("nijg6er")
	m.addCustomBanphrase("nij6ger")
	m.addCustomBanphrase("nij66er")
	m.addCustomBanphrase("nijk6er")
	m.addCustomBanphrase("nijkger")
	m.addCustomBanphrase("nij6ker")
	m.addCustomBanphrase("nijgker")

	m.addCustomBanphrase("n1jgger")
	m.addCustomBanphrase("n1jg6er")
	m.addCustomBanphrase("n1j6ger")
	m.addCustomBanphrase("n1j66er")
	m.addCustomBanphrase("n1jk6er")
	m.addCustomBanphrase("n1jkger")
	m.addCustomBanphrase("n1j6ker")
	m.addCustomBanphrase("n1jgker")

	m.addCustomBanphrase("n!jgger")
	m.addCustomBanphrase("n!jg6er")
	m.addCustomBanphrase("n!j6ger")
	m.addCustomBanphrase("n!j66er")
	m.addCustomBanphrase("n!jk6er")
	m.addCustomBanphrase("n!jkger")
	m.addCustomBanphrase("n!j6ker")
	m.addCustomBanphrase("n!jgker")

	m.addCustomBanphrase("nij1gger")
	m.addCustomBanphrase("nij1g6er")
	m.addCustomBanphrase("nij16ger")
	m.addCustomBanphrase("nij166er")

	m.addCustomBanphrase("nij1k6er")
	m.addCustomBanphrase("nij1kger")
	m.addCustomBanphrase("nij16ker")
	m.addCustomBanphrase("nij1gker")

	m.addCustomBanphrase("n1jk6er")
	m.addCustomBanphrase("n1jkger")
	m.addCustomBanphrase("n1j6ker")
	m.addCustomBanphrase("n1jgker")

	m.addCustomBanphrase("nij1gger")

	m.addCustomBanphrase("ni@@er")
	m.addCustomBanphrase("nigger")
	m.addCustomBanphrase("niggger")
	m.addCustomBanphrase("nigggger")
	m.addCustomBanphrase("niggggger")
	m.addCustomBanphrase("nig6er")
	m.addCustomBanphrase("ni6ger")
	m.addCustomBanphrase("ni6grs")
	m.addCustomBanphrase("ni66er")
	m.addCustomBanphrase("nik6er")
	m.addCustomBanphrase("nikger")
	m.addCustomBanphrase("ni6ker")
	m.addCustomBanphrase("nigker")

	m.addCustomBanphrase("n1gger")
	m.addCustomBanphrase("n1g6er")
	m.addCustomBanphrase("n16ger")
	m.addCustomBanphrase("n166er")
	m.addCustomBanphrase("n1k6er")
	m.addCustomBanphrase("n1kger")
	m.addCustomBanphrase("n16ker")
	m.addCustomBanphrase("n1gker")

	m.addCustomBanphrase("n!gger")
	m.addCustomBanphrase("n!g6er")
	m.addCustomBanphrase("n!6ger")
	m.addCustomBanphrase("n!66er")
	m.addCustomBanphrase("n!k6er")
	m.addCustomBanphrase("n!kger")
	m.addCustomBanphrase("n!6ker")
	m.addCustomBanphrase("n!gker")

	m.addCustomBanphrase("ni1gger")
	m.addCustomBanphrase("ni1g6er")
	m.addCustomBanphrase("ni16ger")
	m.addCustomBanphrase("ni166er")

	m.addCustomBanphrase("ni1k6er")
	m.addCustomBanphrase("ni1kger")
	m.addCustomBanphrase("ni16ker")
	m.addCustomBanphrase("ni1gker")

	m.addCustomBanphrase("n1k6er")
	m.addCustomBanphrase("n1kger")
	m.addCustomBanphrase("n16ker")
	m.addCustomBanphrase("n1gker")

	m.addCustomBanphrase("nek6er")
	m.addCustomBanphrase("nekger")
	m.addCustomBanphrase("ne6ker")
	m.addCustomBanphrase("negker")
	m.addCustomBanphrase("neg6er")
	m.addCustomBanphrase("negger")
	m.addCustomBanphrase("ne6ger")
	m.addCustomBanphrase("negger")

	m.addCustomBanphrase("ni1gg3r")

	m.addCustomBanphrase("nigg3r")
	m.addCustomBanphrase("nig63r")
	m.addCustomBanphrase("ni6g3r")
	m.addCustomBanphrase("ni663r")
	m.addCustomBanphrase("nik63r")
	m.addCustomBanphrase("nikg3r")
	m.addCustomBanphrase("ni6k3r")
	m.addCustomBanphrase("nigk3r")

	m.addCustomBanphrase("n1gg3r")
	m.addCustomBanphrase("n1g63r")
	m.addCustomBanphrase("n16g3r")
	m.addCustomBanphrase("n1663r")
	m.addCustomBanphrase("n1k63r")
	m.addCustomBanphrase("n1kg3r")
	m.addCustomBanphrase("n16k3r")
	m.addCustomBanphrase("n1gk3r")

	m.addCustomBanphrase("n!gg3r")
	m.addCustomBanphrase("n!g63r")
	m.addCustomBanphrase("n!6g3r")
	m.addCustomBanphrase("n!663r")
	m.addCustomBanphrase("n!k63r")
	m.addCustomBanphrase("n!kg3r")
	m.addCustomBanphrase("n!6k3r")
	m.addCustomBanphrase("n!gk3r")

	m.addCustomBanphrase("ni1gg3r")
	m.addCustomBanphrase("ni1g63r")
	m.addCustomBanphrase("ni16g3r")
	m.addCustomBanphrase("ni1663r")

	m.addCustomBanphrase("ni1k63r")
	m.addCustomBanphrase("ni1kg3r")
	m.addCustomBanphrase("ni16k3r")
	m.addCustomBanphrase("ni1gk3r")

	m.addCustomBanphrase("n1k63r")
	m.addCustomBanphrase("n1kg3r")
	m.addCustomBanphrase("n16k3r")
	m.addCustomBanphrase("n1gk3r")

	m.addCustomBanphrase("n3k63r")
	m.addCustomBanphrase("n3kg3r")
	m.addCustomBanphrase("n36k3r")
	m.addCustomBanphrase("n3gk3r")
	m.addCustomBanphrase("n3g63r")
	m.addCustomBanphrase("n3gg3r")
	m.addCustomBanphrase("n36g3r")
	m.addCustomBanphrase("n3gg3r")

	m.addCustomBanphrase("g63r")

	m.banphrases = append(m.banphrases, &filters.Pajbot1Banphrase{
		ID:            -1,
		Name:          "Custom",
		Phrase:        "b00ger",
		Length:        1,
		Operator:      filters.OperatorContains,
		Permanent:     false,
		Warning:       false,
		Notify:        false,
		CaseSensitive: false,
		Enabled:       true,
		SubImmunity:   false,
		RemoveAccents: true,
	})
	err := m.loadPajbot1Banphrases()
	if err != nil {
		return err
	}

	return nil
}

func (m *pajbot1BanphraseFilter) Disable() error {
	return nil
}

func (m *pajbot1BanphraseFilter) Spec() pkg.ModuleSpec {
	return &pajbot1BanphraseSpec
}

func (m *pajbot1BanphraseFilter) BotChannel() pkg.BotChannel {
	return m.botChannel
}

type TimeoutData struct {
	FullMessage string
	Banphrase   pkg.Banphrase
	Username    string
	Channel     string
	Timestamp   time.Time
}

func (m *pajbot1BanphraseFilter) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m *pajbot1BanphraseFilter) check(bot pkg.Sender, source pkg.Channel, text string, action pkg.Action) error {
	originalVariations, lowercaseVariations, err := utils.MakeVariations(text, true)
	if err != nil {
		return err
	}

	for _, bp := range m.banphrases {
		var variations *[]string

		if !bp.IsCaseSensitive() {
			variations = &lowercaseVariations
		} else {
			variations = &originalVariations
		}

		for _, variation := range *variations {
			if bp.Triggers(variation) {
				// fmt.Printf("Banphrase triggered: %#v\n", bp)
				/*
					if bp.IsAdvanced() && source.GetChannel() == "forsen" {
						lol := TimeoutData{
							FullMessage: message.GetText(),
							Banphrase:   bp,
							Username:    user.GetName(),
							Channel:     source.GetChannel(),
							Timestamp:   time.Now().UTC(),
						}
						c := m.server.redis.Get()
						bytes, _ := json.Marshal(&lol)
						c.Do("LPUSH", "pajbot2:timeouts", bytes)
						c.Close()
					}
				*/

				if source.GetChannel() == "krakenbul" || bp.GetID() == -1 {
					reason := fmt.Sprintf("Matched banphrase with name '%s' and id '%d'", bp.GetName(), bp.GetID())
					action.Set(pkg.Timeout{bp.GetLength(), reason})
					action.SetNotifyModerator(bot.MakeUser("pajlada"))
					// fmt.Printf("Banphrase triggered: %#v for user %s", bp, user.GetName())
					return nil
				}
			}

			if !bp.IsAdvanced() {
				if bp.GetID() == -1 {
					fmt.Println("wtf")
				}
				break
			}
		}
	}

	return nil
}

func (m *pajbot1BanphraseFilter) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	if user.IsModerator() || user.IsBroadcaster(source) {
		return nil
	}

	if user.GetName() == "supibot" {
		return nil
	}

	m.check(bot, source, message.GetText(), action)
	m.check(bot, source, user.GetName(), action)

	return nil
}
