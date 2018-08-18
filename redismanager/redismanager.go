package redismanager

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/pkg/common/config"
)

// RedisManager keeps the pool of redis connections
type RedisManager struct {
	Pool *redis.Pool
}

// Init connects to redis and returns redis client
func Init(config config.RedisConfig) (*RedisManager, error) {
	r := &RedisManager{}
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Host)
			if err != nil {
				log.Fatal("An error occured while connecting to redis: ", err)
				return nil, err
			}
			if config.Database >= 0 {
				_, err = c.Do("SELECT", config.Database)
				if err != nil {
					log.Fatal("Error while selecting redis db:", err)
					return nil, err
				}
			}
			return c, err
		},
	}
	r.Pool = pool

	// Ensure that the redis connection works
	conn := r.Pool.Get()
	err := conn.Send("PING")
	if err != nil {
		return nil, err
	}

	return r, nil
}
