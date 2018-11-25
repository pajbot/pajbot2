package utils

import "errors"

// arr = ["a", "b", "c", "d", "e"]
// ChunkStringSlice(arr, 2)
// ret = [ [ "a", "b"], ["c", "d"], ["e"] ]
func ChunkStringSlice(in []string, n int) (out [][]string, err error) {
	if n < 1 {
		err = errors.New("n must be >= 1 in ChunkStringSlice")
		return
	}

	out = make([][]string, CeilInt(float64(len(in))/float64(n)))

	for i := range out {
		out[i] = in[i*n : MinInt((i+1)*n, len(in))]
	}

	return
}
