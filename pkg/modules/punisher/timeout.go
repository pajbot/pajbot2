package punisher

import "time"

type timeout struct {
	received time.Time
	end      time.Time

	duration int
}

func (t *timeout) IsSmaller(o *timeout, leniency int) bool {
	if o.duration == t.duration {
		// We are equal, so we are not smaller XD
		return false
	}

	if o.duration == 0 {
		// Other timeout is a permaban, so we are smaller
		return true
	}

	return o.duration-t.duration > leniency
}

func (t *timeout) Seconds() float64 {
	return time.Until(t.end).Seconds()
}
