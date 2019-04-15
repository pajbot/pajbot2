package config

import "os"

func GetDSN() string {
	const e = `PAJBOT2_DISCORD_BOT_SQL_DSN`
	const defaultValue = `postgres:///pajbot2_discord?sslmode=disable`
	// const defaultValue = `host=/var/run/postgresql database=botsync sslmode=disable`

	if value, ok := os.LookupEnv(e); ok {
		return value
	}

	return defaultValue
}
