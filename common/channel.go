package common

import (
	"database/sql"
	"log"
	"time"

	"github.com/pajlada/pajbot2/sqlmanager"
)

// Channel contains data about the channel
type Channel struct {
	// ID in the database
	ID int

	// Name of the channel, i.e. forsenlol
	Name string

	// Nickname of the channel, i.e. Forsen (could be used as an alias for website)
	Nickname string

	// Channel ID (fetched from the twitch API)
	TwitchChannelID int

	// Access token provided by user on signup
	TwitchAccessToken  string
	TwitchRefreshToken string

	BttvEmotes map[string]Emote // channel and global emotes
	Online     bool
	Uptime     time.Time // time when stream went live, time when last stream ended if not online
}

// FetchAllChannels loads all channels from pb_channel in MySQL
func FetchAllChannels(sql *sqlmanager.SQLManager) ([]Channel, error) {
	var channels []Channel

	rows, err := sql.Session.Query("SELECT id, name, nickname, twitch_channel_id, twitch_access_token, twitch_refresh_token FROM pb_channel")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := Channel{}
		err = c.FetchFromSQL(rows)
		if err != nil {
			// The channel was fetched properly
			channels = append(channels, c)
		}
	}

	return channels, nil
}

// FetchFromSQL populates the given object with data from SQL based on the
// given argument
func (c *Channel) FetchFromSQL(row *sql.Rows) error {
	err := row.Scan(&c.ID, &c.Name, &c.Nickname, &c.TwitchChannelID, &c.TwitchAccessToken, &c.TwitchRefreshToken)
	if err != nil {
		return err
	}
	return nil
}

// InsertToSQL inserts the given channel to SQL
func (c *Channel) InsertToSQL(sql *sqlmanager.SQLManager) error {
	const queryF = `INSERT INTO pb_channel (name) VALUES (?)`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
	}

	// XXX(pajlada): Do we want a SQLChannel middleware to handle null values?

	_, err = stmt.Exec(c.Name)

	if err != nil {
		// XXX
		log.Fatal(err)
	}
	return nil
}
