package pkg

type User interface {
	HasPermission(Permission) bool
	GetName() string
	GetDisplayName() string
	GetID() string
	IsModerator() bool
	IsBroadcaster(Channel) bool
}
