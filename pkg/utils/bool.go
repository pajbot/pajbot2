package utils

// BoolPtr returns a bool pointer of the given bool value
func BoolPtr(v bool) *bool {
	return &v
}

// CheckFlag returns true if the given flag is enabled in the value
func CheckFlag(value uint32, flag uint32) bool {
	return (value & flag) != 0
}
