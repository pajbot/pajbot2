package common

import (
	"database/sql"
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

	// Enabled decides whether we join this channel or not
	Enabled bool

	// Channel ID (fetched from the twitch API)
	TwitchChannelID int64

	// Access token provided by user on signup
	TwitchAccessToken  string
	TwitchRefreshToken string

	BttvEmotes map[string]Emote // channel and global emotes
	Online     bool
	Uptime     time.Time // time when stream went live, time when last stream ended if not online
}

// ChannelSQLWrapper contains data about the channel that's stored in MySQL
type ChannelSQLWrapper struct {
	ID                 int
	Name               string
	Nickname           sql.NullString
	Enabled            int
	TwitchChannelID    sql.NullInt64
	TwitchAccessToken  sql.NullString
	TwitchRefreshToken sql.NullString
}

// FetchAllChannels loads all channels from pb_channel in MySQL
func FetchAllChannels(sql *sqlmanager.SQLManager) ([]Channel, error) {
	var channels []Channel

	rows, err := sql.Session.Query("SELECT id, name, nickname, enabled, twitch_channel_id, twitch_access_token, twitch_refresh_token FROM pb_channel")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := Channel{}
		err = c.FetchFromSQL(rows)
		if err != nil {
			log.Error(err)
		} else {
			// The channel was fetched properly
			channels = append(channels, c)
		}
	}

	return channels, nil
}

/*
FetchFromWrapper is the linker function that links the values of
the ChannelSQLWrapper and the Channel struct together.
This needs to be kept up to date when the table structure of pb_channel
is changed.
*/
func (c *Channel) FetchFromWrapper(w ChannelSQLWrapper) {
	c.ID = w.ID
	c.Name = w.Name
	if w.Nickname.Valid {
		c.Nickname = w.Nickname.String
	}
	if w.Enabled != 0 {
		c.Enabled = true
	}
	if w.TwitchChannelID.Valid {
		c.TwitchChannelID = w.TwitchChannelID.Int64
	}
	if w.TwitchAccessToken.Valid {
		c.TwitchAccessToken = w.TwitchAccessToken.String
	}
	if w.TwitchRefreshToken.Valid {
		c.TwitchRefreshToken = w.TwitchRefreshToken.String
	}
}

// FetchFromSQL populates the given object with data from SQL based on the
// given argument
func (c *Channel) FetchFromSQL(row *sql.Rows) error {
	w := ChannelSQLWrapper{}

	err := row.Scan(&w.ID, &w.Name, &w.Nickname, &w.Enabled, &w.TwitchChannelID, &w.TwitchAccessToken, &w.TwitchRefreshToken)

	if err != nil {
		log.Error(err)
		return err
	}

	c.FetchFromWrapper(w)

	return nil
}

// InsertNewToSQL inserts the given channel to SQL
func (c *Channel) InsertNewToSQL(sql *sqlmanager.SQLManager) error {
	const queryF = `INSERT INTO pb_channel (name) VALUES (?)`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(c.Name)

	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}
	return nil
}

// SQLSetEnabled updates the enabled state of the given channel
func (c *Channel) SQLSetEnabled(sql *sqlmanager.SQLManager, enabled int) error {
	const queryF = `UPDATE pb_channel SET enabled=? WHERE id=?`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(enabled, c.ID)

	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}
	return nil
}
