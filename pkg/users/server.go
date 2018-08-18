package users

import (
	"database/sql"
)

type server struct {
	sql *sql.DB
}

var _server server

func InitServer(_sql *sql.DB) error {
	_server.sql = _sql

	return nil
}
