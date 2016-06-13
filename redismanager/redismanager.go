package redismanager

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis" // are you ok with this package?
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/helper"
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

func (r *Redismanager) SetPoints(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("ZADD", channel+":users:points", user.Points, user.Name)
	conn.Flush()
}

func (r *Redismanager) IncrPoints(channel string, user *common.User, incrby int) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("ZINCRBY", channel+":users:points", incrby, user.Name)
	conn.Flush()
}

func (r *Redismanager) newUser(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("HSET", channel+":users:lastseen", user.Name, time.Now().Unix())
	conn.Send("ZADD", channel+":users:points", user.Points, user.Name)
	conn.Send("HSET", channel+":users:level", user.Name, float64(r.getLevel(0.1, user)))
	conn.Flush()
}

func (r *Redismanager) SetLevel(channel string, user *common.User, level int) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("HSET", channel+":users:level", user.Name, float64(level)+0.2)
	conn.Flush()
}

func (r *Redismanager) ResetLevel(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("HSET", channel+":users:level", user.Name, float64(r.getLevel(0.1, user))+0.1)
	conn.Flush()
}

func (r *Redismanager) UpdateUser(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	if user.Name == channel {
		r.SetLevel(channel, user, 1500)
	}
	conn.Send("HSET", channel+":users:lastseen", user.Name, time.Now().Unix())
	conn.Flush()
}

// GetUser fills out missing fields of the given User object
// and creates new user in redis if the user doesnt exist
func (r *Redismanager) GetUser(channel string, user *common.User) {
	fmt.Println(user.Name)
	conn := r.Pool.Get()
	defer conn.Close()
	exist, err := conn.Do("HEXISTS", channel+":users:lastseen", user.Name)
	e, _ := redis.Bool(exist, err)
	if e {
		conn.Send("HGET", channel+":users:level", user.Name)
		conn.Send("ZSCORE", channel+":users:points", user.Name)
		conn.Send("HGET", channel+":users:lastseen", user.Name)
		conn.Flush()
		// can this be done in a loop somehow?
		// Level
		res, err := conn.Receive()
		level, _ := redis.Float64(res, err)
		user.Level = r.getLevel(level, user)
		// Points
		res, err = conn.Receive()
		user.Points, _ = redis.Int(res, err)
		// LastSeen
		res, err = conn.Receive()
		lastseen, _ := redis.String(res, err)
		user.LastSeen, _ = time.Parse(time.UnixDate, lastseen)
	} else {
		r.newUser(channel, user)
		r.GetUser(channel, user)
	}
}

/*
.1 : level not set manually, the bot will change it automatically

.2 : level set manually, the bot will not change it until !level reset
global level is always set manually
*/
func (r *Redismanager) getLevel(channel float64, user *common.User) int {

	if float64(user.Level) > channel {
		return user.Level
	}
	status := helper.Round(channel, 1) * 10
	if status == 2 {
		return int(channel)
	}
	if user.ChannelOwner {
		return 1500
	} else if user.Mod {
		return 500
	} else if user.Sub {
		return 250
	} else {
		return 100
	}
}
