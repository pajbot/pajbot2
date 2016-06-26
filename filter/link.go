package filter

import "github.com/mvdan/xurls"

// LinkFilter xD
func LinkFilter(m string) []string {
	return xurls.Relaxed.FindAllString(m, -1)
}
