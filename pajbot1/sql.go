package pajbot1

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"
	"unicode"

	twitch "github.com/gempir/go-twitch-irc"
	normalize "github.com/pajlada/lidl-normalize"
	"github.com/pajlada/pajbot2/bots"
	"github.com/pajlada/pajbot2/common/config"
)

type Pajbot1 struct {
	Session *sql.DB

	EnabledBanphrases []Pajbot1Banphrase
}

type OperatorType int

const (
	OperatorContains OperatorType = iota
	OperatorStartsWith
	OperatorEndsWith
	OperatorExact
)

type Pajbot1Banphrase struct {
	ID     int
	Name   string
	Phrase string
	Length int

	// "contains" or "startswith" or "endswith" or "exact"
	Operator      OperatorType // handled
	Permanent     bool
	Warning       bool
	Notify        bool
	CaseSensitive bool // handled
	Enabled       bool // handled
	SubImmunity   bool
	RemoveAccents bool // handled, and a little bit more
}

// Init creates an instance of the SQL Manager
func Init(config config.Pajbot1Config) *Pajbot1 {
	m := &Pajbot1{}

	db, err := sql.Open("mysql", config.SQL.DSN)
	if err != nil {
		log.Fatal("Error connecting to MySQL:", err)
	}
	// TODO: Close database

	m.Session = db

	return m
}

func (p *Pajbot1) LoadBanphrases() error {
	const queryF = `SELECT * FROM tb_banphrase`

	stmt, err := p.Session.Prepare(queryF)
	if err != nil {
		return err
	}

	rows, err := stmt.Query()
	if err != nil {
		return err
	}

	p.EnabledBanphrases = []Pajbot1Banphrase{}

	var bp Pajbot1Banphrase
	var operatorString string

	for rows.Next() {
		err = rows.Scan(&bp.ID, &bp.Name, &bp.Phrase, &bp.Length, &bp.Permanent, &bp.Warning, &bp.Notify, &bp.CaseSensitive, &bp.Enabled, &operatorString, &bp.SubImmunity, &bp.RemoveAccents)
		if err != nil {
			return err
		}

		if bp.Enabled {
			if operatorString == "contains" {
				bp.Operator = OperatorContains
			} else if operatorString == "startswith" {
				bp.Operator = OperatorStartsWith
			} else if operatorString == "endswith" {
				bp.Operator = OperatorEndsWith
			} else if operatorString == "exact" {
				bp.Operator = OperatorExact
			}
			p.EnabledBanphrases = append(p.EnabledBanphrases, bp)
		}
	}

	return nil
}

func HandleContains(phrase, text string) bool {
	return strings.Contains(text, phrase)
}

func HandleExact(phrase, text string) bool {
	return phrase == text
}

func HandleStartsWith(phrase, text string) bool {
	return strings.HasPrefix(text, phrase)
}

func HandleEndsWith(phrase, text string) bool {
	return strings.HasSuffix(text, phrase)
}

type TimeoutData struct {
	FullMessage string
	Banphrase   Pajbot1Banphrase
	Username    string
	Channel     string
	Timestamp   time.Time
}

func DoTimeout(bot *bots.TwitchBot, bp *Pajbot1Banphrase, channel string, user twitch.User, message *bots.TwitchMessage) {
	if !bp.RemoveAccents {
		return
	}

	lol := TimeoutData{
		FullMessage: message.Text,
		Banphrase:   *bp,
		Username:    user.Username,
		Channel:     channel,
		Timestamp:   time.Now().UTC(),
	}
	c := bot.Redis.Pool.Get()
	bytes, _ := json.Marshal(&lol)
	c.Do("LPUSH", "pajbot2:timeouts", bytes)
	c.Close()
	// bot.Timeout(channel, user, bp.Length, "Matched banphrase with name \""+bp.Name+"\"")
	// bot.Say(channel, user.Username+" matched banphrase with name "+bp.Name)
	// log.Println("Matched banphrase with name \"" + bp.Name + "\"")
}

func lowercaseAll(in []string) []string {
	out := make([]string, len(in))

	for i, v := range in {
		out[i] = strings.ToLower(v)
	}

	return out
}

type removeFunc func(rune) bool

func removeInStringFunc(in string, predicate removeFunc) string {
	outBytes := make([]rune, len(in))
	length := 0
	for _, r := range in {
		if !predicate(r) {
			outBytes[length] = r
			length++
		}
	}

	return string(outBytes[:length])
}

const LatinCapitalLetterBegin = 0x41
const LatinCapitalLetterEnd = 0x5A

const LatinSmallLetterBegin = 0x61
const LatinSmallLetterEnd = 0x7A

func isNotLatinLetter(r rune) bool {
	return !((r >= LatinSmallLetterBegin && r <= LatinSmallLetterEnd) || (r >= LatinCapitalLetterBegin && r <= LatinCapitalLetterEnd))
}

func insertUnique(text string, target *[]string) {
	for _, v := range *target {
		if v == text {
			// log.Printf("Not inserting because %s == %s\n", v, text)
			return
		}
	}

	// log.Println("Inserting", text)

	*target = append(*target, text)
}

func (p *Pajbot1) CheckBanphrases(next bots.Handler) bots.Handler {
	return bots.HandlerFunc(func(bot *bots.TwitchBot, channel string, user twitch.User, message *bots.TwitchMessage) {
		normalizedMessage, err := normalize.Normalize(message.Text)
		if err != nil {
			log.Println("Err:", err)
			return
		}

		originalVariations := []string{
			// Full message
			message.Text,
		}

		// Full message with all spaces removed
		insertUnique(removeInStringFunc(message.Text, unicode.IsSpace), &originalVariations)

		// Full message with all spaces and non-latin letters removed
		insertUnique(removeInStringFunc(message.Text, isNotLatinLetter), &originalVariations)

		// Normalized message
		insertUnique(normalizedMessage, &originalVariations)

		// Normalized message with all spaces removed
		insertUnique(removeInStringFunc(normalizedMessage, unicode.IsSpace), &originalVariations)

		// Normalized message with all spaces non-latin letters removed
		insertUnique(removeInStringFunc(normalizedMessage, isNotLatinLetter), &originalVariations)

		lowercaseVariations := lowercaseAll(originalVariations)

		if user.UserType != "" && user.Username != "pajlada" {
			next.HandleMessage(bot, channel, user, message)
			return
		}

		variationPrinted := make([]bool, len(originalVariations))
		for _, bp := range p.EnabledBanphrases {
			var phrase string
			var variations *[]string

			if !bp.CaseSensitive {
				phrase = strings.ToLower(bp.Phrase)
				variations = &lowercaseVariations
			} else {
				phrase = bp.Phrase
				variations = &originalVariations
			}

			for variationIndex, variation := range *variations {
				if !variationPrinted[variationIndex] && user.Username == "pajlada" {
					log.Printf("Testing %d: '%s'\n", variationIndex, variation)
					variationPrinted[variationIndex] = true
				}
				switch bp.Operator {
				case OperatorContains:
					if HandleContains(phrase, variation) {
						DoTimeout(bot, &bp, channel, user, message)
						return
					}

				case OperatorExact:
					if HandleExact(phrase, variation) {
						DoTimeout(bot, &bp, channel, user, message)
						return
					}

				case OperatorStartsWith:
					if HandleStartsWith(phrase, variation) {
						DoTimeout(bot, &bp, channel, user, message)
						return
					}

				case OperatorEndsWith:
					if HandleEndsWith(phrase, variation) {
						DoTimeout(bot, &bp, channel, user, message)
						return
					}
				}

				if !bp.RemoveAccents {
					break
				}
			}
		}

		next.HandleMessage(bot, channel, user, message)
	})
}
