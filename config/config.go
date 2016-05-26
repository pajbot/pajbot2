package config

/*
The Config contains all the data required to connect
to the twitch IRC servers
*/
type Config struct {
	Pass       string
	Nick       string
	BrokerPort string

	RedisPw string
	RedisIP string

	TLSKey  string
	TLSCert string

	ToWeb   chan map[string]interface{}
	FromWeb chan map[string]interface{}
}

/*
GetConfig returns a singleton instance? of the config object
*/
func GetConfig() Config {
	config := &Config{
		Pass:    "oauth:ai1b8xkefpjek6gckutjovus82nulx",
		Nick:    "testaccount_420",
		RedisIP: ":6379",

		ToWeb:   make(chan map[string]interface{}),
		FromWeb: make(chan map[string]interface{}),
	}
	return *config
}
