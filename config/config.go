package config

type Config struct {
	Pass       string
	Nick       string
	BrokerPort string

	RedisPw string
	RedisIP string

	TlsKey  string
	TlsCert string

	ToWeb   chan map[string]interface{}
	FromWeb chan map[string]interface{}
}

func GetConfig() Config {
	config := &Config{
		Pass: "oauth:xD",
		Nick: "nuulsbot",

		RedisIP: ":6379",

		ToWeb:   make(chan map[string]interface{}),
		FromWeb: make(chan map[string]interface{}),
	}
	return *config
}
