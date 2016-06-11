package redismanager

import (
	"fmt"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/common"
)

type Redismanager struct {
	Pool *redis.Pool
}

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

func (r *Redismanager) UpdateUser(channel string, user *common.User) {
	conn := r.Pool.Get()
	defer conn.Close()
	conn.Send("ZADD", channel+":points", user.Points, user.Name)
	//                ^^^^^^^^^^^^^^^^^^  is this how you do it? LUL
	conn.Send("HSET", channel+":lastseen", user.Name, time.Now().Unix())
	conn.Flush()
	log.Printf("saved new user %s\n", user.Name)
}

// GetUser takes the channel and pointer to user object and fills out missing fields
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
