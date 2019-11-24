package web

import "regexp"

const (
	// ErrInvalidUserName is returned if an invalid username was posted to a router that expected a valid username
	ErrInvalidUserName = "invalid username"
)

var (
	singleUserName = regexp.MustCompile(`\w+`)
	userNameList   = regexp.MustCompile(`[\w\,]+`)
)

func isValidUserName(input string) bool {
	return singleUserName.FindString(input) == input
}
