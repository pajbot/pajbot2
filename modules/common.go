package modules

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pajlada/pajbot2/bot"
	"github.com/pajlada/pajbot2/common/basemodule"
	"github.com/pajlada/pajbot2/plog"
)

var log = plog.GetLogger()

/*
Module state json:
{
	"enabled": true(enabled), false(disabled), null(default value),
	"level_required": 350(level 350 required), null(default level required (usually 100)),
	"level_bypass": 350(level 350 required to ignore module check), null(default level required to bypass(usually -1(unbypassable)))
	"settings": {
		"roulette_amount": null, // unset (default value)
		"setting2": 5, // int value set
		"setting3": [1, 2, 3], // list value set
		"setting4": {
			"int": 5,
		},
		"setting5": {
			"string": "xD",
		}
	}
}
*/

type moduleSettings struct {
	BanphraseDefaultLength int64 `json:"banphrase_default_length"`
}

type rouletteModuleSettings struct {
	MaxRouletteAmount int64 `json:"max_roulette_amount"`
}

type duelModuleSettings struct {
	MaxDuelAmount int64 `json:"max_duel_amount"`
}

func getModuleSettingsBytes(b *bot.Bot, m *basemodule.BaseModule) []byte {
	conn := b.Redis.Pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("HGET", b.Channel.Name+":modules:settings", m.ID))
	if err != nil {
		return nil
	}
	return data
}

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
