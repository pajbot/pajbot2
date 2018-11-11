package utils

func StringContains(needle string, haystack []string) bool {
	for _, straw := range haystack {
		if needle == straw {
			return true
		}
	}

	return false
}

func SBKey(m map[string]bool) (keys []string) {
	keys = make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}

	return
}
