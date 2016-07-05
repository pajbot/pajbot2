package apirequest

import "time"

type Stream struct {
	ID         string
	Online     bool
	Created    time.Time
	Game       string
	IsPlaylist bool
	Viewers    int
}

type Channel struct {
	ID        string
	Status    string
	Game      string
	UpdatedAt time.Time
	Views     int
	Followers int
	Partner   bool
}
