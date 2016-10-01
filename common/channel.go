package common

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pajlada/pajbot2/sqlmanager"
)

const channelQ = "SELECT id, name, nickname, enabled, bot_id FROM pb_channel"

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

	// XXX: this should probably we renamed to BotAcountID instead of naming it BotID or bot_id everywhere
	BotID int

	Emotes ExtensionEmotes

	Online bool
	Uptime time.Time // time when stream went live, time when last stream ended if not online
}

// ChannelSQLWrapper contains data about the channel that's stored in MySQL
type ChannelSQLWrapper struct {
	ID              int
	Name            string
	Nickname        sql.NullString
	Enabled         int
	TwitchChannelID sql.NullInt64
	BotID           int
}

// FetchAllChannels loads all channels from pb_channel in MySQL
func FetchAllChannels(sql *sqlmanager.SQLManager, botID int) ([]Channel, error) {
	var channels []Channel

	const queryF = channelQ + " WHERE bot_id=?"

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(botID)
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
	c.BotID = w.BotID
}

// FetchFromSQL populates the given object with data from SQL based on the
// given argument
func (c *Channel) FetchFromSQL(row *sql.Rows) error {
	w := ChannelSQLWrapper{}

	err := row.Scan(&w.ID, &w.Name, &w.Nickname, &w.Enabled, &w.BotID)

	if err != nil {
		log.Error(err)
		return err
	}

	c.FetchFromWrapper(w)

	return nil
}

// FetchFromSQLRow populates the given object with data from SQL based on the
// given argument
func (c *Channel) FetchFromSQLRow(row *sql.Row) error {
	w := ChannelSQLWrapper{}

	err := row.Scan(&w.ID, &w.Name, &w.Nickname, &w.Enabled, &w.BotID)

	if err != nil {
		return err
	}

	c.FetchFromWrapper(w)

	return nil
}

// InsertNewToSQL inserts the given channel to SQL
func (c *Channel) InsertNewToSQL(sql *sqlmanager.SQLManager) error {
	const queryF = `INSERT INTO pb_channel (name, bot_id) VALUES (?, ?)`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(c.Name, c.BotID)

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

// SQLSetBotID updates the enabled state of the given channel
func (c *Channel) SQLSetBotID(sql *sqlmanager.SQLManager, botID int) error {
	const queryF = `UPDATE pb_channel SET bot_id=? WHERE id=?`

	stmt, err := sql.Session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}

	_, err = stmt.Exec(botID, c.ID)

	if err != nil {
		// XXX
		log.Fatal(err)
		return err
	}
	return nil
}

// GetChannel xD
func GetChannel(session *sql.DB, name string) (Channel, error) {
	const queryF = channelQ + " WHERE name=?"

	stmt, err := session.Prepare(queryF)
	if err != nil {
		return Channel{}, err
	}

	var c Channel

	err = c.FetchFromSQLRow(stmt.QueryRow(name))

	switch {
	case err == sql.ErrNoRows:
		log.Error(err)
		return Channel{}, fmt.Errorf("No channel with the name %s", name)

	case err != nil:
		log.Error(err)
		return Channel{}, err
	}

	return c, nil
}
