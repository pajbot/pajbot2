package pkg

type Account interface {
	// full id of user (i.e. 11148817)
	ID() string

	// full lowercase name (i.e. pajlada)
	Name() string
}
