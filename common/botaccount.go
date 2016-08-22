package common

import "database/sql"

const botAccountQ = "SELECT id, name, twitch_access_token, twitch_refresh_token FROM pb_bot_account"

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
	const queryF = botAccountQ + " WHERE name=?"

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
		return BotAccount{}, err

	case err != nil:
		return BotAccount{}, err
	}

	ba := BotAccount{
		ID:   outID,
		Name: outName,
		TwitchCredentials: TwitchClientCredentials{
			AccessToken:  outTwitchAccessToken,
			RefreshToken: outTwitchRefreshToken,
		},
	}

	return ba, nil
}

// GetAllBotAccounts xD
func GetAllBotAccounts(session *sql.DB) ([]BotAccount, error) {
	const queryF = botAccountQ
	var botAccounts []BotAccount

	var outID int
	var outName string
	var outTwitchAccessToken string
	var outTwitchRefreshToken string

	rows, err := session.Query(queryF)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&outID, &outName, &outTwitchAccessToken, &outTwitchRefreshToken)
		if err != nil {
			log.Error(err)
		} else {
			ba := BotAccount{
				ID:   outID,
				Name: outName,
				TwitchCredentials: TwitchClientCredentials{
					AccessToken:  outTwitchAccessToken,
					RefreshToken: outTwitchRefreshToken,
				},
			}
			botAccounts = append(botAccounts, ba)
		}
	}

	return botAccounts, nil
}
