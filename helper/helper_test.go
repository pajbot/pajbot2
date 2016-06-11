package helper

import (
	"fmt"
	"testing"
)

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
	fmt.Println("xD")
}
