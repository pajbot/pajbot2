package modules

import (
	"database/sql"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/pkg/common/config"
	"github.com/pajlada/pajbot2/pkg/pubsub"
)

type server struct {
	redis      *redis.Pool
	sql        *sql.DB
	oldSession *sql.DB
	pubSub     *pubsub.PubSub
}

var _server server

func InitServer(redis *redis.Pool, _sql *sql.DB, pajbot1Config config.Pajbot1Config, pubSub *pubsub.PubSub) error {
	var err error

	_server.redis = redis
	_server.sql = _sql
	_server.oldSession, err = sql.Open("mysql", pajbot1Config.SQL.DSN)
	_server.pubSub = pubSub
	if err != nil {
		return err
	}

	return nil
}
