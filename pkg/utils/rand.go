package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

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
