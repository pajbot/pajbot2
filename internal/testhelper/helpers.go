package testhelper

import (
	"testing"
)

func AssertStringsEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("failed asserting that \"%s\" is expected \"%s\"", actual, expected)
	}
}

func AssertIntsEqual(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("failed asserting that \"%d\" is expected \"%d\"", actual, expected)
	}
}

func AssertInt32sEqual(t *testing.T, expected, actual int32) {
	if expected != actual {
		t.Errorf("failed asserting that \"%d\" is expected \"%d\"", actual, expected)
	}
}

func AssertTrue(t *testing.T, actual bool, errorMessage string) {
	if !actual {
		t.Error(errorMessage)
	}
}

func AssertFalse(t *testing.T, actual bool, errorMessage string) {
	if actual {
		t.Error(errorMessage)
	}
}

func AssertNil(t *testing.T, actual interface{}, errorMessage string) {
	if actual != nil {
		t.Error(errorMessage)
	}
}

func AssertNotNil(t *testing.T, actual interface{}, errorMessage string) {
	if actual == nil {
		t.Error(errorMessage)
	}
}

func AssertInterfacesEqual(t *testing.T, expected, actual interface{}, errorMessage string) {
	if actual != expected {
		t.Error(errorMessage)
	}
}

func AssertStringSlicesEqual(t *testing.T, expected, actual []string) {
	if actual == nil {
		t.Errorf("actual slice was nil")
	}

	if len(actual) != len(expected) {
		t.Errorf("actual slice was not the same length as expected slice")
	}

	for i, v := range actual {
		if v != expected[i] {
			t.Errorf("actual slice value \"%s\" was not equal to expected value \"%s\" at index \"%d\"", v, expected[i], i)
		}
	}
}

func AssertErrorsEqual(t *testing.T, expected, actual error) {
	if expected != actual {
		t.Errorf("failed asserting that error \"%s\" is expected \"%s\"", actual, expected)
	}
}
