package config

import "os"

var (
	WebStaticPath = stringEnv("PAJBOT2_WEB_PATH", "../../web/")
)

func stringEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
