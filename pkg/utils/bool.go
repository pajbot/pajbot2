package utils

// GetTrueP returns a bool true pointer
func GetTrueP() *bool {
	b := true
	return &b
}

// GetFalseP returns a bool false pointer
func GetFalseP() *bool {
	b := false
	return &b
}

// CheckFlag returns true if the given flag is enabled in the value
func CheckFlag(value uint32, flag uint32) bool {
	return (value & flag) != 0
}
