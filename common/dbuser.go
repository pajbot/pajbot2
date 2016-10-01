package common

import "database/sql"

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
	const queryF = `INSERT INTO pb_user(name, type, twitch_access_token, twitch_refresh_token) VALUES (?, ?, ?, ?)`

	stmt, err := session.Prepare(queryF)
	if err != nil {
		log.Debugf("error: %s", err)
		return err
	}
	_, err = stmt.Exec(name, userType, accessToken, refreshToken)
	if err != nil {
		log.Debugf("error: %s", err)
		return err
	}

	return nil
}

// GetDBUser returns a user by name and type
func GetDBUser(session *sql.DB, name string, userType string) (DBUser, error) {
	const queryF = userQ + " WHERE name=? AND type=?"

	stmt, err := session.Prepare(queryF)
	if err != nil {
		return DBUser{}, err
	}
	var outID int
	var outName string
	var outTwitchAccessToken string
	var outTwitchRefreshToken string

	err = stmt.QueryRow(name).Scan(&outID, &outName, &outTwitchAccessToken, &outTwitchRefreshToken)
	switch {
	case err == sql.ErrNoRows:
		return DBUser{}, err

	case err != nil:
		return DBUser{}, err
	}

	user := DBUser{
		ID:   outID,
		Name: outName,
		Type: userType,
		TwitchCredentials: TwitchClientCredentials{
			AccessToken:  outTwitchAccessToken,
			RefreshToken: outTwitchRefreshToken,
		},
	}

	return user, nil
}

// GetDBUsersByType returns a list of DB Users
func GetDBUsersByType(session *sql.DB, userType string) ([]DBUser, error) {
	const queryF = userQ + " WHERE type=?"
	var dbUsers []DBUser

	var outID int
	var outName string
	var outTwitchAccessToken string
	var outTwitchRefreshToken string

	stmt, err := session.Prepare(queryF)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userType)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&outID, &outName, &outTwitchAccessToken, &outTwitchRefreshToken)
		if err != nil {
			log.Error(err)
		} else {
			ba := DBUser{
				ID:   outID,
				Name: outName,
				Type: "bot",
				TwitchCredentials: TwitchClientCredentials{
					AccessToken:  outTwitchAccessToken,
					RefreshToken: outTwitchRefreshToken,
				},
			}
			dbUsers = append(dbUsers, ba)
		}
	}

	return dbUsers, nil
}
