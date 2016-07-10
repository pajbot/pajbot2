package helper

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strings"
)

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

// SplitUint64 spits a uint64 value into two uint32 values
func SplitUint64(val uint64) (uint32, uint32) {
	return uint32(val >> 32), uint32(val)
}

// CombineUint32 returns the packed uint64 value of two uint32's
func CombineUint32(val1 uint32, val2 uint32) uint64 {
	r := uint64(val1)
	r = r << 32
	r += uint64(val2)
	return r
}

// CheckFlag returns true if the given flag is enabled in the value
func CheckFlag(value uint32, flag uint32) bool {
	return (value & flag) != 0
}

// NewStringPtr returns the pointer to the given string
func NewStringPtr(s string) *string {
	return &s
}

// GetTriggers returns a list of strings that have been parsed in accordance
// to the command rules
func GetTriggers(message string) []string {
	return strings.Split(strings.Replace(strings.ToLower(message), "!", "", 1), " ")
}

// GetTriggersN returns a list of strings that have been parsed in accordance
// to the command rules. Offset by N
func GetTriggersN(message string, n int) []string {
	triggers := GetTriggers(message)
	if len(triggers) >= n {
		return triggers[n:]
	}
	return []string{}
}

// GetTriggersKC returns a list of strings that have been parsed in accordance
// to the command rules, but keeps the case
func GetTriggersKC(message string) []string {
	return strings.Split(strings.Replace(message, "!", "", 1), " ")
}

// RemoveNewlines replaces all \r and \n with spaces
func RemoveNewlines(s string) string {
	s = strings.Replace(s, "\r", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	return s
}

/*
RandIntN generates a number between min and max (inclusive).
Only works with positive numbers.
*/
func RandIntN(min int, max int) (int, error) {
	if min < 0 {
		return 0, fmt.Errorf("min must be a positive number")
	}

	if min > max {
		return 0, fmt.Errorf("min must be bigger than max")
	}

	toN := big.NewInt(int64(max - min))

	val, err := rand.Int(rand.Reader, toN)

	if err != nil {
		return 0, err
	}

	return int(val.Int64()), nil
}
