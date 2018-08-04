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

func Uint32ToBytes(value uint32) (bytes []byte) {
	bytes = make([]byte, 4)

	bytes[0] = (byte)((value >> 24) & 0xFF)
	bytes[1] = (byte)((value >> 16) & 0xFF)
	bytes[2] = (byte)((value >> 8) & 0xFF)
	bytes[3] = (byte)((value) & 0xFF)

	return bytes
}

func Uint64ToBytes(value uint64) (bytes []byte) {
	bytes = make([]byte, 8)

	bytes[0] = (byte)((value >> 56) & 0xFF)
	bytes[1] = (byte)((value >> 48) & 0xFF)
	bytes[2] = (byte)((value >> 40) & 0xFF)
	bytes[3] = (byte)((value >> 32) & 0xFF)
	bytes[4] = (byte)((value >> 24) & 0xFF)
	bytes[5] = (byte)((value >> 16) & 0xFF)
	bytes[6] = (byte)((value >> 8) & 0xFF)
	bytes[7] = (byte)((value) & 0xFF)

	return bytes
}
