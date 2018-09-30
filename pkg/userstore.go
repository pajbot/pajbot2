package pkg

type UserStore interface {
	// Input: Lowercased twitch usernames
	// Returns: user IDs as strings in no specific order, and a bool indicating whether the user needs to exhaust the list first and wait
	GetIDs([]string) map[string]string

	GetID(string) string

	GetName(string) string

	// Input: list of twitch IDs
	// Returns: map of twitch IDs pointing at twitch usernames
	GetNames([]string) map[string]string
}
