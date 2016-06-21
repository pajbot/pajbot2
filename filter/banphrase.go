package filter

// BanAction xD , just an idea, this wont really work xD
type BanAction struct {
	/*
	   1 not important
	   3 timeout worthy
	   5 sure timeout
	   7 long timeout
	   10 perm ban
	   1-3 will be sent to dashboard
	*/
	Level   int
	Reason  string   // matched ascii filter, matched link filter etc..
	Matches []string // matched parts of msg
	Matched bool
}
