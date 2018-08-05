package pkg

type User interface {
	HasGlobalPermission(Permission) bool
	HasChannelPermission(Channel, Permission) bool
	GetName() string
	GetDisplayName() string
	GetID() string
	IsModerator() bool
	IsBroadcaster(Channel) bool
}
