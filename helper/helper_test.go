package helper

import "testing"

func TestSum(t *testing.T) {
	var sumTests = []struct {
		input    []int
		expected int
	}{
		{
			input:    []int{1, 2, 3},
			expected: 6,
		},
		{
			input:    []int{-5, 5},
			expected: 0,
		},
		{
			input:    []int{5, 5, 5, 5},
			expected: 20,
		},
	}

	for _, tt := range sumTests {
		res := Sum(tt.input)

		if res != tt.expected {
			t.Errorf("%d is not equal to %d", res, tt.expected)
		}
	}
}

func TestRound(t *testing.T) {
	var roundTests = []struct {
		inputVal    float64
		inputPlaces int
		expected    float64
	}{
		{
			inputVal:    69.69696969,
			inputPlaces: 2,
			expected:    69.70,
		},
		{
			inputVal:    123.4005,
			inputPlaces: 1,
			expected:    123.4,
		},
	}

	for _, tt := range roundTests {
		res := Round(tt.inputVal, tt.inputPlaces)

		if res != tt.expected {
			t.Errorf("%f is not equal to %f", res, tt.expected)
		}
	}
}
