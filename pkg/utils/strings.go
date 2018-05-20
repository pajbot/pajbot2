package utils

import (
	"strings"
	"unicode"

	normalize "github.com/pajlada/lidl-normalize"
)

const latinCapitalLetterBegin = 0x41
const latinCapitalLetterEnd = 0x5A

const latinSmallLetterBegin = 0x61
const latinSmallLetterEnd = 0x7A

func isNotLatinLetter(r rune) bool {
	return !((r >= latinSmallLetterBegin && r <= latinSmallLetterEnd) || (r >= latinCapitalLetterBegin && r <= latinCapitalLetterEnd))
}

type removeFunc func(rune) bool

func removeInStringFunc(in string, predicate removeFunc) string {
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
	InsertUnique(removeInStringFunc(text, unicode.IsSpace), &originalVariations)

	// Full message with all spaces and non-latin letters removed
	InsertUnique(removeInStringFunc(text, isNotLatinLetter), &originalVariations)

	if doNormalize {
		normalizedMessage, err := normalize.Normalize(text)
		if err != nil {
			return nil, nil, err
		}
		// Normalized message
		InsertUnique(normalizedMessage, &originalVariations)

		// Normalized message with all spaces removed
		InsertUnique(removeInStringFunc(normalizedMessage, unicode.IsSpace), &originalVariations)

		// Normalized message with all spaces non-latin letters removed
		InsertUnique(removeInStringFunc(normalizedMessage, isNotLatinLetter), &originalVariations)
	}

	return originalVariations, lowercaseAll(originalVariations), nil
}
