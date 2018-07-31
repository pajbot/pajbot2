package pkg

type Sender interface {
	Say(Channel, string)
	Mention(Channel, User, string)
	Timeout(Channel, User, int, string)
	GetPoints(Channel, User) uint64

	// give or remove points from user in channel
	EditPoints(Channel, User, int32) uint64
	BulkEdit(string, []string, int32)

	GetUserStore() UserStore
}
