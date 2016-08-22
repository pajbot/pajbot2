package boss

import (
	"github.com/pajlada/pajbot2/common"
	"github.com/pajlada/pajbot2/common/config"
	"github.com/pajlada/pajbot2/pbtwitter"
	"github.com/pajlada/pajbot2/redismanager"
	"github.com/pajlada/pajbot2/sqlmanager"
)

// Boss is the struct that keeps all of the IRC Connections in place
// One IRC Connection is equivalent to one "bot account"
type Boss struct {
	IRCConnections map[string]*Irc
	Redis          *redismanager.RedisManager
	SQL            *sqlmanager.SQLManager
	Twitter        *pbtwitter.Client
}

// Init intializes a boss struct
func Init(config *config.Config) Boss {
	boss := Boss{
		IRCConnections: make(map[string]*Irc),
		Redis:          redismanager.Init(config),
		SQL:            sqlmanager.Init(config),
	}
	boss.Twitter = pbtwitter.Init(config, boss.Redis)

	// Shared config between every IRC Instance
	c := IRCConfig{
		BrokerHost: *config.BrokerHost,
		BrokerPass: *config.BrokerPass,
		Redis:      boss.Redis,
		SQL:        boss.SQL,
		Twitter:    boss.Twitter,
		Quit:       config.Quit,
		Silent:     config.Silent,
	}

	botAccounts, err := common.GetAllBotAccounts(boss.SQL.Session)
	if err != nil {
		log.Error(err)
	}

	for _, botAccount := range botAccounts {
		c.BrokerLogin = config.BrokerLogin + "-" + botAccount.Name
		log.Debugf("Starting bot %s", botAccount.Name)
		ircConnection := InitIRCConnection(c, botAccount)
		boss.IRCConnections[botAccount.Name] = ircConnection
	}

	return boss
}
