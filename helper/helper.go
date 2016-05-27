package helper

/*
Sum returns the sum of the given slice of ints
*/
func Sum(s []int) int {
	var x int
	for i := range s {
		x += s[i]
	}
	return x
}
