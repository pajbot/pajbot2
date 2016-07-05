package redismanager

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// LoadTwitterFollows loads followed users
func (r *RedisManager) LoadTwitterFollows() ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	// a list doesnt seem right but what would be better? json? pajaW
	users, err := redis.Strings(conn.Do("LRANGE", "twitterfollows", "0", "10000000"))
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("expired")
	}
	return users, nil
}

// SaveTwitterFollows saves followed users and sets an expire for 2 hours
func (r *RedisManager) SaveTwitterFollows(users []string) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("DEL", "twitterfollows")
	for _, user := range users {
		conn.Send("LPUSH", "twitterfollows", user)
	}
	conn.Send("EXPIRE", "twitterfollows", 2*60*60) // 2 hours
	conn.Flush()
}
