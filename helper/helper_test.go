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

func TestSplitUint64(t *testing.T) {
	var splitTests = []struct {
		input     uint64
		expectedA uint32
		expectedB uint32
	}{
		{
			input:     4294967298,
			expectedA: 1,
			expectedB: 2,
		},
		{
			input:     30064771074,
			expectedA: 7,
			expectedB: 2,
		},
	}

	for _, tt := range splitTests {
		resA, resB := SplitUint64(tt.input)

		if resA != tt.expectedA {
			t.Errorf("%d is not equal to %d", resA, tt.expectedA)
		}

		if resB != tt.expectedB {
			t.Errorf("%d is not equal to %d", resB, tt.expectedA)
		}
	}
}

func TestCheckFlag(t *testing.T) {
	var tests = []struct {
		inputValue uint32
		inputFlag  uint32
		expected   bool
	}{
		{
			inputValue: 3,
			inputFlag:  1,
			expected:   true,
		},
		{
			inputValue: 2,
			inputFlag:  1,
			expected:   false,
		},
		{
			inputValue: 16,
			inputFlag:  1,
			expected:   false,
		},
		{
			inputValue: 18,
			inputFlag:  2,
			expected:   true,
		},
		{
			inputValue: 17,
			inputFlag:  1,
			expected:   true,
		},
	}

	for _, tt := range tests {
		res := CheckFlag(tt.inputValue, tt.inputFlag)

		if res != tt.expected {
			t.Errorf("%t is not equal to %t (%d - %d)", res, tt.expected, tt.inputValue, tt.inputFlag)
		}
	}
}
