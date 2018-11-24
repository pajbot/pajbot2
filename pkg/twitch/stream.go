package twitch

import "github.com/pajlada/pajbot2/pkg"

var _ pkg.Stream = &Stream{}

type Stream struct {
}

func (s *Stream) Live() bool {
	return false
}
