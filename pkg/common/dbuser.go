package common

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
