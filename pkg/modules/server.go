package modules

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/sqlmanager"
)

type server struct {
	redis      *redis.Pool
	sql        *sqlmanager.SQLManager
	oldSession *sql.DB
}

var _server server

func InitServer(redis *redis.Pool, _sql *sqlmanager.SQLManager, pajbot1Config config.Pajbot1Config) error {
	var err error

	_server.redis = redis
	_server.sql = _sql
	_server.oldSession, err = sql.Open("mysql", pajbot1Config.SQL.DSN)
	if err != nil {
		return err
	}

	return nil
}
