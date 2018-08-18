package utils

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

var packingTests = []struct {
	partA  uint32
	partB  uint32
	packed uint64
}{
	{
		partA:  1,
		partB:  2,
		packed: 4294967298,
	},
	{
		partA:  1,
		partB:  4,
		packed: 4294967300,
	},
	{
		partA:  7,
		partB:  2,
		packed: 30064771074,
	},
}

func TestSplitUint64(t *testing.T) {
	for _, tt := range packingTests {
		resA, resB := SplitUint64(tt.packed)

		if resA != tt.partA {
			t.Errorf("%d is not equal to %d", resA, tt.partA)
		}

		if resB != tt.partB {
			t.Errorf("%d is not equal to %d", resB, tt.partB)
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

func TestCombineUint32(t *testing.T) {
	for _, tt := range packingTests {
		res := CombineUint32(tt.partA, tt.partB)

		if res != tt.packed {
			t.Errorf("%d is not equal to %d", res, tt.packed)
		}
	}
}
