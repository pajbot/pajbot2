package modules

import (
	"database/sql"

	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/pajlada/pajbot2/sqlmanager"
)

type server struct {
	redis      *redismanager.RedisManager
	sql        *sqlmanager.SQLManager
	oldSession *sql.DB
}

var _server server

func InitServer(redis *redismanager.RedisManager, _sql *sqlmanager.SQLManager, pajbot1Config config.Pajbot1Config) error {
	var err error

	_server.redis = redis
	_server.sql = _sql
	_server.oldSession, err = sql.Open("mysql", pajbot1Config.SQL.DSN)
	if err != nil {
		return err
	}

	return nil
}
