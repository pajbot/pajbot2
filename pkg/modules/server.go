package modules

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/pkg/common/config"
)

type server struct {
	redis      *redis.Pool
	sql        *sql.DB
	oldSession *sql.DB
}

var _server server

func InitServer(redis *redis.Pool, _sql *sql.DB, pajbot1Config config.Pajbot1Config) error {
	var err error

	_server.redis = redis
	_server.sql = _sql
	_server.oldSession, err = sql.Open("mysql", pajbot1Config.SQL.DSN)
	if err != nil {
		return err
	}

	return nil
}
