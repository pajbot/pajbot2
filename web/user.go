package web

import (
	"strings"
	"time"
)

func getUserPayload(channel, username string) interface{} {
	username = strings.ToLower(username)
	if !redis.IsValidUser(channel, username) {
		return newError("user not found")
	}
	u := redis.LoadUser(channel, username)
	return user{
		Name:                username,
		DisplayName:         u.DisplayName,
		Points:              int64(u.Points),
		Level:               int64(u.Level),
		TotalMessageCount:   int64(u.TotalMessageCount),
		OfflineMessageCount: int64(u.OfflineMessageCount),
		OnlineMessageCount:  int64(u.OnlineMessageCount),
		LastSeen:            u.LastSeen.Format(time.UnixDate),
		LastActive:          u.LastActive.Format(time.UnixDate),
	}
}
