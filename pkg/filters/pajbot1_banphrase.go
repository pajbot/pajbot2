package filters

import (
	"database/sql"
	"regexp"
	"strings"
)

// BanphraseOperator is a banphrase operator
type BanphraseOperator int

const (
	OperatorContains BanphraseOperator = iota
	OperatorStartsWith
	OperatorEndsWith
	OperatorExact
	OperatorRegex
)

// Pajbot1Banphrase is a banphrase loaded from the old pajbot1 database
type Pajbot1Banphrase struct {
	ID     int
	Name   string
	Phrase string
	Length int

	// "contains" or "startswith" or "endswith" or "exact"
	Operator      BanphraseOperator // handled
	Permanent     bool
	Warning       bool
	Notify        bool
	CaseSensitive bool // handled
	Enabled       bool // handled
	SubImmunity   bool
	RemoveAccents bool // handled, and a little bit more
}

func handleContains(phrase, text string) bool {
	return strings.Contains(text, phrase)
}

func handleExact(phrase, text string) bool {
	return phrase == text
}

func handleStartsWith(phrase, text string) bool {
	return strings.HasPrefix(text, phrase)
}

func handleEndsWith(phrase, text string) bool {
	return strings.HasSuffix(text, phrase)
}

func handleRegex(phrase, text string) bool {
	re, err := regexp.Compile(phrase)
	if err != nil {
		return false
	}

	return re.MatchString(text)
}

func (f *Pajbot1Banphrase) Triggers(text string) bool {
	// log.Println("Do we", f.Phrase, "trigger", text, "? forsenThink")
	switch f.Operator {
	case OperatorContains:
		if handleContains(f.Phrase, text) {
			return true
		}

	case OperatorExact:
		if handleExact(f.Phrase, text) {
			return true
		}

	case OperatorStartsWith:
		if handleStartsWith(f.Phrase, text) {
			return true
		}

	case OperatorEndsWith:
		if handleEndsWith(f.Phrase, text) {
			return true
		}
	case OperatorRegex:
		if handleRegex(f.Phrase, text) {
			return true
		}
	}

	return false
}

func (f *Pajbot1Banphrase) IsCaseSensitive() bool {
	return f.CaseSensitive
}

func (f *Pajbot1Banphrase) IsAdvanced() bool {
	return f.RemoveAccents
}

func (f *Pajbot1Banphrase) GetLength() int {
	return f.Length
}

func (f *Pajbot1Banphrase) GetName() string {
	return f.Name
}

func (f *Pajbot1Banphrase) GetID() int {
	return f.ID
}

func (f *Pajbot1Banphrase) LoadScan(rows *sql.Rows) error {
	var operatorString string
	err := rows.Scan(&f.ID, &f.Name, &f.Phrase, &f.Length, &f.Permanent, &f.Warning, &f.Notify, &f.CaseSensitive, &f.Enabled, &operatorString, &f.SubImmunity, &f.RemoveAccents)
	if err != nil {
		return err
	}

	if !f.CaseSensitive {
		f.Phrase = strings.ToLower(f.Phrase)
	}

	if operatorString == "contains" {
		f.Operator = OperatorContains
	} else if operatorString == "startswith" {
		f.Operator = OperatorStartsWith
	} else if operatorString == "endswith" {
		f.Operator = OperatorEndsWith
	} else if operatorString == "exact" {
		f.Operator = OperatorExact
	} else if operatorString == "regex" {
		f.Operator = OperatorRegex
	}

	return nil
}
