package pkg

type Channel interface {
	Say(string, string)
	Timeout(string, User, int, string)
}
