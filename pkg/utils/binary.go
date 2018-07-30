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

func Int32ToBytes(value int32) (bytes []byte) {
	bytes = make([]byte, 4)

	bytes[0] = (byte)((value >> 24) & 0xFF)
	bytes[1] = (byte)((value >> 16) & 0xFF)
	bytes[2] = (byte)((value >> 8) & 0xFF)
	bytes[3] = (byte)((value) & 0xFF)

	return bytes
}
