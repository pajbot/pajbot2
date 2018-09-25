package pkg

type Permission uint64

const (
	PermissionNone       Permission = 0
	PermissionReport     Permission = 1 << 0
	PermissionRaffle     Permission = 1 << 1
	PermissionAdmin      Permission = 1 << 2
	PermissionModeration Permission = 1 << 3
)

// GetPermissionBit converts a string (i.e. "admin") to the binary value it represents.
// 0b100 in this example
func GetPermissionBit(s string) Permission {
	if s == "report" {
		return PermissionReport
	}
	if s == "raffle" {
		return PermissionRaffle
	}
	if s == "admin" {
		return PermissionAdmin
	}
	if s == "moderation" {
		return PermissionModeration
	}

	return PermissionNone
}

// GetPermissionBits converts a list of strings (i.e. ["admin", "raffle"]) to the binary value they represent.
// 0b110 in this example
func GetPermissionBits(permissionNames []string) (permissions Permission) {
	for _, permissionName := range permissionNames {
		permission := GetPermissionBit(permissionName)
		permissions |= permission
	}

	return
}
