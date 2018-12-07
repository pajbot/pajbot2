package ws

// MessageType is used to help redirect a message to the proper connections
type MessageType uint8

// All available MessageTypes
const (
	MessageTypeAll MessageType = iota
	MessageTypeNone
	MessageTypeCLR
	MessageTypeDashboard
)

// WSMessage xD
type WSMessage struct {
	Channel string

	MessageType MessageType

	// LevelRequired <=0 means the message does not require authentication, otherwise
	// authentication is required and the users level must be equal to or above
	// the LevelRequired value
	LevelRequired int

	Payload *Payload
}
