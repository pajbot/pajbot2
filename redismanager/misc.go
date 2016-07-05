package redismanager

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

/*
should we do the same with bttv emotes or store everything in an apicache hash instead?
hashes dont support EXPIRE so we'd have to do it manually
*/

// LoadTwitterFollows loads followed users
func (r *RedisManager) LoadTwitterFollows() ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	bs, err := redis.Bytes(conn.Do("GET", "twitterfollows"))
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("expired")
	}
	var users []string
	err = json.Unmarshal(bs, &users)
	return users, err
}

// SaveTwitterFollows saves followed users and sets an expire for 2 hours
func (r *RedisManager) SaveTwitterFollows(users []string) {
	conn := r.Pool.Get()
	defer conn.Close()
	bs, err := json.Marshal(users)
	if err != nil {
		log.Error(err)
	}
	conn.Send("DEL", "twitterfollows")
	conn.Send("SET", "twitterfollows", bs)
	conn.Send("EXPIRE", "twitterfollows", 2*60*60) // 2 hours
	conn.Flush()
}
