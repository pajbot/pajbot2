package redismanager

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis" // are you ok with this package?
	"github.com/pajlada/pajbot2/common"
)

type Redismanager struct {
	Pool *redis.Pool
}

// Init connects to redis and returns redis client
func Init(config *common.Config) *Redismanager {
	r := &Redismanager{}
	pool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", config.RedisHost)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return c, err
	}, 69)
	// forsenGASM
	r.Pool = pool
	log.Println("connected to redis")
	return r
}

// only globally banned users and admins have a level in global redis
func (r *Redismanager) UpdateGlobalUser(channel string, user *common.User, u *common.GlobalUser) {
	log.Printf("redis: user: %s  channel: %s", user.Name, channel)
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("HSET", "global:lastactive", user.Name, time.Now().Unix())
	conn.Send("HSET", "global:channel", user.Name, channel)
	conn.Flush()
}

func (r *Redismanager) GetGlobalUser(channel string, user *common.User, u *common.GlobalUser) {
	conn := r.Pool.Get()
	defer conn.Close()
	exist, err := conn.Do("HEXISTS", "global:lastactive", user.Name)
	e, _ := redis.Bool(exist, err)
	if e {
		conn.Send("HGET", "global:level", user.Name)
		conn.Send("HGET", "global:lastactive", user.Name)
		conn.Send("HGET", "global:channel", user.Name)
		conn.Flush()
		// can this be done in a loop somehow?
		// Level
		res, err := conn.Receive()
		level, _ := redis.Int(res, err) // will be 0 unless user is admin or globally banned
		if level > user.Level {
			user.Level = level
		}
		// LastActive
		res, err = conn.Receive()
		t, _ := redis.String(res, err)
		u.LastActive, _ = time.Parse(time.UnixDate, t)
		// Channel
		res, err = conn.Receive()
		u.Channel, _ = redis.String(res, err)
	} else {
		r.UpdateGlobalUser(channel, user, u)
		r.GetGlobalUser(channel, user, u)
	}
	r.UpdateGlobalUser(channel, user, u)
}

func (r *Redismanager) UpdateUser(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("ZADD", channel+":points", user.Points, user.Name)
	//                ^^^^^^^^^^^^^^^^^^  is this how you do it? LUL
	conn.Send("HSET", channel+":lastseen", user.Name, time.Now().Unix())
	conn.Flush()
	log.Printf("saved new user %s\n", user.Name)
}

// GetUser fills out missing fields of the given User object
// and creates new user in redis if the user doesnt exist
func (r *Redismanager) GetUser(channel string, user *common.User) {
	fmt.Println(user.Name)
	conn := r.Pool.Get()
	defer conn.Close()
	exist, err := conn.Do("HEXISTS", channel+":lastseen", user.Name)
	e, _ := redis.Bool(exist, err)
	if e {
		conn.Send("ZSCORE", channel+":points", user.Name)
		conn.Send("HGET", channel+":lastseen", user.Name)
		conn.Flush()
		// can this be done in a loop somehow?
		res, err := conn.Receive()
		user.Points, _ = redis.Int(res, err)
		res, err = conn.Receive()
		lastseen, _ := redis.String(res, err)
		user.LastSeen, _ = time.Parse(time.UnixDate, lastseen)
	} else {
		r.UpdateUser(channel, user)
		r.GetUser(channel, user)
	}
}
