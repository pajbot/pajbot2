package pkg

type User interface {
	// Has channel or global permission
	HasPermission(Channel, Permission) bool

	// Has global permission
	HasGlobalPermission(Permission) bool

	// Has channel permission
	HasChannelPermission(Channel, Permission) bool

	GetName() string
	GetDisplayName() string
	GetID() string
	IsModerator() bool
	IsBroadcaster() bool
	IsVIP() bool
	GetBadges() map[string]int
}
