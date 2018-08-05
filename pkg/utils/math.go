package utils

func Abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func Abs64(n int64) int64 {
	if n < 0 {
		return -n
	}

	return n
}
