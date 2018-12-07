package web

import "regexp"

// funny haHAA
const (
	ErrInvalidUserName = "https://i.imgur.com/r7FGMh8.png"
)

var (
	singleUserName = regexp.MustCompile(`\w+`)
	userNameList   = regexp.MustCompile(`[\w\,]+`)
	rawURL         = regexp.MustCompile(`[\w\,\/]+`)
)

func isValidUserName(input string) bool {
	if singleUserName.FindString(input) != input {
		return false
	}
	return true
}

func isValidURL(url string) bool {
	if rawURL.FindString(url) != url {
		return false
	}
	return true
}
