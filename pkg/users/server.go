package users

import "github.com/pajlada/pajbot2/sqlmanager"

type server struct {
	sql *sqlmanager.SQLManager
}

var _server server

func InitServer(_sql *sqlmanager.SQLManager) error {
	_server.sql = _sql

	return nil
}
