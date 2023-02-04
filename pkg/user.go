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
	IsSubscriber() bool
	GetBadges() map[string]int

	// Set the ID of this user
	// Will return an error if the ID of this user was already set
	SetID(string) error

	// Set the Name of this user
	// Will return an error if the Name of this user was already set
	SetName(string) error
}
