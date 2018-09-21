package common

import (
	"database/sql"
	"log"
)

// DBUser xD
type DBUser struct {
	// ID in the database
	ID int

	// Name of the user, i.e. snusbot
	Name string

	// Type of the user
	Type string

	TwitchCredentials TwitchClientCredentials
}

const userQ = "SELECT id, name, twitch_access_token, twitch_refresh_token FROM pb_user"

// CreateDBUser creates a bot in the pb_bot table
func CreateDBUser(session *sql.DB, name string, accessToken string, refreshToken string, userType string) error {
	const queryF = `INSERT INTO pb_bot(name, twitch_access_token, twitch_refresh_token) VALUES (?, ?, ?, ?)`

	_, err := session.Exec(queryF, name, userType, accessToken, refreshToken)
	if err != nil {
		log.Printf("error: %s", err)
		return err
	}

	return nil
}

// CreateBot xD
func CreateBot(session *sql.DB, name string, accessToken string, refreshToken string) error {
	const queryF = `INSERT INTO pb_bot(name, twitch_access_token, twitch_refresh_token) VALUES (?, ?, ?)`

	_, err := session.Exec(queryF, name, accessToken, refreshToken)
	if err != nil {
		log.Printf("error: %s", err)
		return err
	}

	return nil
}
