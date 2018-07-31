package pkg

type UserStore interface {
	// Input: Lowercased twitch usernames
	// Returns: user IDs as strings in no specific order, and a bool indicating whether the user needs to exhaust the list first and wait
	GetIDs([]string) map[string]string
}
