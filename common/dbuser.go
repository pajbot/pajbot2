package common

import "database/sql"

// CreateDBUser creates a bot in the pb_bot table
func CreateDBUser(session *sql.DB, name string, accessToken string, refreshToken string) {
	const queryF = `INSERT INTO pb_user(name, twitch_access_token, twitch_refresh_token) VALUES (?, ?, ?)`

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
