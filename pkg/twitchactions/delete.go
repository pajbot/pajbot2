package twitchactions

type deleteAction struct {
	message string
}

func (m *deleteAction) Message() string {
	return m.message
}
