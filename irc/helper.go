package irc

/*
Sum returns the sum of the given slice of ints

TODO: Should this be in just helper/helper.go instead?
*/
func Sum(s []int) int {
	var x int
	for i := range s {
		x += s[i]
	}
	return x
}
