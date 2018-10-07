package pkg

type UserContext interface {
	GetContext(channelID, userID string) []string

	// The message must be preformatted. i.e. [2018-09-13 502:501:201] username: message
	AddContext(channelID, userID, message string)
}
