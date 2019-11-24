package web

import "regexp"

// funny haHAA
const (
	ErrInvalidUserName = "https://i.imgur.com/r7FGMh8.png"
)

var (
	singleUserName = regexp.MustCompile(`\w+`)
	userNameList   = regexp.MustCompile(`[\w\,]+`)
)

func isValidUserName(input string) bool {
	return singleUserName.FindString(input) == input
}
