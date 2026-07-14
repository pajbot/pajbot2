package utils

import (
	"strconv"
	"time"
)

func ParseTwitchDuration(t string, defaultUnit time.Duration, defaultTime time.Duration) time.Duration {
	if seconds, err := strconv.Atoi(t); err == nil {
		// No suffix = treat as seconds
		return time.Duration(time.Duration(seconds) * defaultUnit)
	}

	d, err := time.ParseDuration(t)
	if err != nil {
		return defaultTime
	}
	return d
}
