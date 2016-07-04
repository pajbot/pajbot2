package common

import "time"

// User xD
type User struct {
	ID                  int
	Name                string
	DisplayName         string
	Mod                 bool
	Sub                 bool
	Turbo               bool
	ChannelOwner        bool
	Type                string // admin , staff etc
	Level               int
	OnlineMessageCount  int
	OfflineMessageCount int
	Points              int
	LastSeen            time.Time // should this be time.Time or int/float?
	LastActive          time.Time
}

const noPing = string("\u05C4")

// NameNoPing xD
func (u *User) NameNoPing() string {
	return string(u.DisplayName[0]) + noPing + u.DisplayName[1:]
}
