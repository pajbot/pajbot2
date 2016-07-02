package filter

import "strings"

// ContainsWord xD
func ContainsWord(msg string, bannedWords []BannedWord) int {
	var lvl int
	m := strings.ToLower(msg)
	for _, word := range bannedWords {
		if strings.Contains(m, word.Word) {
			if word.Level > lvl {
				lvl = word.Level
			}
		}
	}
	return lvl
}

func ContainsLink(links []string, bannedLinks []BannedLink) int {
	var lvl int
	for _, link := range links {
		link = strings.ToLower(link)
		for _, bannedLink := range bannedLinks {
			if strings.Contains(link, bannedLink.Link) || strings.Contains(bannedLink.Link, link) {
				if bannedLink.Level > lvl {
					lvl = bannedLink.Level
				}
			}
		}
	}
	return lvl
}
