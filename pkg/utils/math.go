package utils

import "math"

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

/*
Sum returns the sum of the given slice of ints
*/
func Sum(s []int) int {
	var x int
	for i := range s {
		x += s[i]
	}
	return x
}

func FloorInt(f float64) int {
	return int(math.Floor(f))
}

func CeilInt(f float64) int {
	return int(math.Ceil(f))
}

// Round returns the rounded value of a float64 up to N places
func Round(val float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}

	return b
}
