package command

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

// response types
const (
	Say     ResponseTypeEnum = "say"
	Whisper ResponseTypeEnum = "whisper"
	Reply   ResponseTypeEnum = "reply"
)

// enabled types
const (
	Yes         EnabledEnum = "yes"
	No          EnabledEnum = "no"
	OnlineOnly  EnabledEnum = "online_only"
	OfflineOnly EnabledEnum = "offline_only"
)

// ResponseTypeEnum xD
type ResponseTypeEnum string

// Scan xD
func (u *ResponseTypeEnum) Scan(value interface{}) error {
	*u = ResponseTypeEnum(string(value.([]uint8)))
	return nil
}

// Value xD
func (u *ResponseTypeEnum) Value() (driver.Value, error) {
	return string(*u), nil
}

// EnabledEnum xD
type EnabledEnum string

// SQLCommand xD
type SQLCommand struct {
	ID           int
	ChannelID    int
	Triggers     string
	Response     string
	ResponseType ResponseTypeEnum
	Level        int
	CooldownAll  int
	CooldownUser int
	Enabled      EnabledEnum
	CostPoints   int
	Filters      string // TODO: this is a SET('banphrases', 'linkchecker')
}

// ReadSQLCommand reads a single command
func ReadSQLCommand(rows *sql.Rows) Command {
	sqlCommand := &SQLCommand{}
	err := rows.Scan(&sqlCommand.ID, &sqlCommand.ChannelID, &sqlCommand.Triggers,
		&sqlCommand.Response, &sqlCommand.ResponseType)
	if err != nil {
		log.Error(err)
		return nil
	}
	c := TextCommand{
		BaseCommand: BaseCommand{
			ID:       sqlCommand.ID,
			Triggers: strings.Split(sqlCommand.Triggers, "|"),
		},
		Response: sqlCommand.Response,
	}
	return &c
}

// Insert calls insert on the given sql session
func (command *SQLCommand) Insert(session *sql.DB) int64 {
	const queryF = `INSERT INTO pb_command(channel_id, triggers, response) VALUES (?, ?, ?)`

	stmt, err := session.Prepare(queryF)
	if err != nil {
		// XXX
		log.Fatal(err)
	}
	res, err := stmt.Exec(command.ChannelID, command.Triggers, command.Response)
	if err != nil {
		// XXX
		log.Fatal(err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		// XXX
		log.Fatal(err)
	}
	log.Debugf("Added new command with ID %d", lastID)
	return lastID
}

// Delete deletes the command from the database
func (command *SQLCommand) Delete(session *sql.DB) error {
	// Ensure that the command ID is set
	if command.ID == 0 {
		return errors.New("Invalid SQLCommand used in Delete")
	}

	const queryF = `DELETE FROM pb_command WHERE id=?`

	stmt, err := session.Prepare(queryF)
	if err != nil {
		return err
	}
	res, err := stmt.Exec(command.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Make sure that exactly one row was deleted
	if rowsAffected != 1 {
		return fmt.Errorf("Rows affected is %d when it should be 1", rowsAffected)
	}
	return nil
}
