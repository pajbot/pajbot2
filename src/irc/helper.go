package irc

func Sum(s []int) int {
	var x int
	for i := range s {
		x += s[i]
	}
	return x
}
