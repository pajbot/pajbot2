package modules

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/filters"
	"github.com/pajlada/pajbot2/pkg/utils"
)

type Pajbot1BanphraseFilter struct {
	server *server

	banphrases []pkg.Banphrase
}

func NewPajbot1BanphraseFilter() *Pajbot1BanphraseFilter {
	return &Pajbot1BanphraseFilter{
		server: &_server,
	}
}

func (m *Pajbot1BanphraseFilter) loadPajbot1Banphrases() error {
	const queryF = `SELECT * FROM tb_banphrase`

	session := m.server.oldSession

	stmt, err := session.Prepare(queryF)
	if err != nil {
		return err
	}

	rows, err := stmt.Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		var bp filters.Pajbot1Banphrase
		err = bp.LoadScan(rows)
		if err != nil {
			return err
		}

		if bp.Enabled {
			m.banphrases = append(m.banphrases, &bp)
		}
	}

	return nil
}

func (m *Pajbot1BanphraseFilter) Register() error {
	err := m.loadPajbot1Banphrases()
	if err != nil {
		return err
	}

	return nil
}

func (m Pajbot1BanphraseFilter) Name() string {
	return "Pajbot1BanphraseFilter"
}

type TimeoutData struct {
	FullMessage string
	Banphrase   pkg.Banphrase
	Username    string
	Channel     string
	Timestamp   time.Time
}

func (m Pajbot1BanphraseFilter) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	return nil
}

func (m Pajbot1BanphraseFilter) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	originalVariations, lowercaseVariations, err := utils.MakeVariations(message.GetText(), true)
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
				if bp.IsAdvanced() && source.GetChannel() == "forsen" {
					lol := TimeoutData{
						FullMessage: message.GetText(),
						Banphrase:   bp,
						Username:    user.GetName(),
						Channel:     source.GetChannel(),
						Timestamp:   time.Now().UTC(),
					}
					c := m.server.redis.Pool.Get()
					bytes, _ := json.Marshal(&lol)
					c.Do("LPUSH", "pajbot2:timeouts", bytes)
					c.Close()
				}

				if source.GetChannel() == "krakenbul" && !user.IsModerator() {
					reason := fmt.Sprintf("Matched banphrase with name '%s' and id '%d'", bp.GetName(), bp.GetID())
					action.SetTimeout(bp.GetLength(), reason)
				}
				return nil
			}

			if !bp.IsAdvanced() {
				break
			}
		}
	}

	return nil
}
