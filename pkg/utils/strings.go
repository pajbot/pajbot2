package utils

import (
	"bufio"
	"strconv"
	"strings"
	"unicode"

	normalize "github.com/pajlada/lidl-normalize"
)

const latinCapitalLetterBegin = 0x41
const latinCapitalLetterEnd = 0x5A

const latinSmallLetterBegin = 0x61
const latinSmallLetterEnd = 0x7A

const numbersBegin = 0x30
const numbersEnd = 0x39

func IsNotLatinLetter(r rune) bool {
	return !((r >= latinSmallLetterBegin && r <= latinSmallLetterEnd) || (r >= latinCapitalLetterBegin && r <= latinCapitalLetterEnd))
}

func IsNotAlphanumeric(r rune) bool {
	if r >= latinSmallLetterBegin && r <= latinSmallLetterEnd {
		return false
	}
	if r >= latinCapitalLetterBegin && r <= latinCapitalLetterEnd {
		return false
	}

	if r >= numbersBegin && r <= numbersEnd {
		return false
	}

	return true
}

type RemoveFunc func(rune) bool

func RemoveInStringFunc(in string, predicate RemoveFunc) string {
	outBytes := make([]rune, len(in))
	length := 0
	for _, r := range in {
		if !predicate(r) {
			outBytes[length] = r
			length++
		}
	}

	return string(outBytes[:length])
}

// InsertUnique inserts a string into a target string splice if it doesn't already exist in it
func InsertUnique(text string, target *[]string) {
	for _, v := range *target {
		if v == text {
			return
		}
	}

	*target = append(*target, text)
}

func lowercaseAll(in []string) []string {
	out := make([]string, len(in))

	for i, v := range in {
		out[i] = strings.ToLower(v)
	}

	return out
}

// MakeVariations makes normal-case and lowercase variatinos of a string
func MakeVariations(text string, doNormalize bool) ([]string, []string, error) {
	originalVariations := []string{
		// Full message
		text,
	}

	// Full message with all spaces removed
	InsertUnique(RemoveInStringFunc(text, unicode.IsSpace), &originalVariations)

	// Full message with all spaces and non-latin letters removed
	InsertUnique(RemoveInStringFunc(text, IsNotLatinLetter), &originalVariations)

	// Message with all non-alpha-numberic letters removed
	InsertUnique(RemoveInStringFunc(text, IsNotAlphanumeric), &originalVariations)

	if doNormalize {
		normalizedMessage, err := normalize.Normalize(text)
		if err != nil {
			return nil, nil, err
		}
		// Normalized message
		InsertUnique(normalizedMessage, &originalVariations)

		// Normalized message with all spaces removed
		InsertUnique(RemoveInStringFunc(normalizedMessage, unicode.IsSpace), &originalVariations)

		// Normalized message with all spaces non-latin letters removed
		InsertUnique(RemoveInStringFunc(normalizedMessage, IsNotLatinLetter), &originalVariations)

		// Normalized message with all spaces non-alpha-numberic letters removed
		InsertUnique(RemoveInStringFunc(normalizedMessage, IsNotAlphanumeric), &originalVariations)
	}

	return originalVariations, lowercaseAll(originalVariations), nil
}

// IsValidUsername verifies that the given user name is not empty, and that it only contains valid characters (0-9 a-z A-Z _)
// It does not to any length-checks, as there are no unbroken rules as to how short or long a user name can be
func IsValidUsername(username string) bool {
	if username == "" {
		return false
	}

	for _, r := range username {
		// Numbers || uppercase a-z || lowercase a-z || underscore
		if (r >= 0x30 && r <= 0x39) || (r >= 0x41 && r <= 0x5A) || (r >= 0x61 && r <= 0x7A) || r == 0x5F {
			continue
		}

		return false
	}

	return true
}

// IsValidUserID verifies that the given user id is not empty, and that it only contains valid characters
// It does not to any length-checks, as there are no clear rules as to how short or long a user ID can be
func IsValidUserID(userID string) bool {
	if userID == "" {
		return false
	}

	for _, r := range userID {
		// Numbers
		if r >= 0x30 && r <= 0x39 {
			continue
		}

		return false
	}

	return true
}

func FilterChannelName(username string) string {
	username = strings.TrimPrefix(username, "#")

	if IsValidUsername(username) {
		return username
	}

	return ""
}

func FilterUsername(username string) string {
	username = strings.TrimPrefix(username, "@")
	username = strings.TrimSuffix(username, ",")
	username = strings.TrimSuffix(username, ":")

	if IsValidUsername(username) {
		return username
	}

	return ""
}

func FilterUsernames(potentialUsernames []string) (usernames []string) {
	for _, s := range potentialUsernames {
		if IsValidUsername(s) {
			usernames = append(usernames, strings.ToLower(s))
		}
	}

	return
}

func FilterUserIDs(potentialUserIDs []string) (userIDs []string) {
	for _, s := range potentialUserIDs {
		if IsValidUserID(s) {
			userIDs = append(userIDs, s)
		}
	}

	return
}

const noPing = string("\u05C4")

func UnpingifyUsername(username string) string {
	return string(username[0]) + noPing + username[1:]
}

// ReadArg reads a string until \n and trims all whitespace
func ReadArg(reader *bufio.Reader) string {
	untrimmed, _ := reader.ReadString('\n')

	return strings.TrimSpace(untrimmed)
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

// NewStringPtr returns the pointer to the given string
func NewStringPtr(s string) *string {
	return &s
}

// FloatToString converts a float value to a string with sane defaults
func Float32ToString(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}
