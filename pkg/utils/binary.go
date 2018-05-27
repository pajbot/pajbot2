package utils

// NOTE: bytes must be 8 length here
func BytesToUint64(bytes []uint8) (v uint64) {
	if len(bytes) != 8 {
		panic("Invalid bytes array length sent to BytesToUint64")
	}

	for i, b := range bytes {
		v += uint64(b) << uint(7-i)
	}
	return
}
