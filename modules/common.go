package modules

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/plog"
)

var log = plog.GetLogger()

func isModuleEnabled(b *bot.Bot, id string, defValue bool) bool {
	conn := b.Redis.Pool.Get()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", b.Channel.Name+":modules:"+id+":state"))
	if err != nil {
		log.Errorf("An error occured: %s", err)
		return false
	}

	if !exists {
		// The user has not set any information about this module
		return defValue
	}

	return true
}
