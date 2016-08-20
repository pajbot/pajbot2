package common

import (
	"database/sql"
	"fmt"
)

// BotAccount xD
type BotAccount struct {
	// ID in the database
	ID int

	// Name of the bot, i.e. snusbot
	Name string

	TwitchCredentials TwitchClientCredentials
}

// CreateBotAccount creates a bot in the pb_bot table
func CreateBotAccount(session *sql.DB, name string, accessToken string, refreshToken string) {
	const queryF = `INSERT INTO pb_bot_account(name, twitch_access_token, twitch_refresh_token) VALUES (?, ?, ?)`

	stmt, err := session.Prepare(queryF)
	if err != nil {
		log.Debugf("error: %s", err)
		return
	}
	_, err = stmt.Exec(name, accessToken, refreshToken)
	if err != nil {
		log.Debugf("error: %s", err)
		return
	}
}

// GetBotAccount xD
func GetBotAccount(session *sql.DB, name string) (BotAccount, error) {
	const queryF = `SELECT id, name, twitch_access_token, twitch_refresh_token FROM pb_bot_account WHERE name=?`

	stmt, err := session.Prepare(queryF)
	if err != nil {
		log.Debugf("error: %s", err)
		return BotAccount{}, err
	}

	var outID int
	var outName string
	var outTwitchAccessToken string
	var outTwitchRefreshToken string

	err = stmt.QueryRow(name).Scan(&outID, &outName, &outTwitchAccessToken, &outTwitchRefreshToken)
	switch {
	case err == sql.ErrNoRows:
		return BotAccount{}, fmt.Errorf("No bot account with the name %s", name)

	case err != nil:
		return BotAccount{}, err
	}

	ba := BotAccount{
		ID:   outID,
		Name: outName,
		TwitchCredentials: TwitchClientCredentials{
			TwitchAccessToken:  outTwitchAccessToken,
			TwitchRefreshToken: outTwitchRefreshToken,
		},
	}

	return ba, nil
}
