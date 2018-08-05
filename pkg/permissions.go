package pkg

type Permission uint64

const (
	PermissionNone   Permission = 0
	PermissionReport Permission = 1 << 0
	PermissionRaffle Permission = 1 << 1
	PermissionAdmin  Permission = 1 << 2
)

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

	return PermissionNone
}
