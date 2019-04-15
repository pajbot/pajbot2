package main

type Action struct {
	// Valid types:
	// - unmute
	Type string

	GuildID string
	UserID  string
	RoleID  string
}
