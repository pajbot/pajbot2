package sqlmanager

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL Driver
	"github.com/pajlada/pajbot2/pkg/common/config"
)

// SQLManager keeps a pool of sql connections or some shit like that
type SQLManager struct {
	Session *sql.DB
}

// Init creates an instance of the SQL Manager
func Init(config config.SQLConfig) *SQLManager {
	m := &SQLManager{}

	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal("Error connecting to MySQL:", err)
	}
	// TODO: Close database

	m.Session = db

	return m
}
