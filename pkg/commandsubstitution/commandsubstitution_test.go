package commandsubstitution

import "testing"

var (
	args = map[string]Substitution{
		"user": &User{
			name:  "Pajlada",
			level: 100,
		},
	}
)

func TestShouldMatch(t *testing.T) {
	const message = `a $(user.name) b`
	const expected = `a Pajlada b`

	result, err := Substitute(message, args)
	assertErrorsEqual(t, nil, err)
	assertStringsEqual(t, expected, result)
}

func TestFilter(t *testing.T) {
	const message = `a $(user.name|toupper) b`
	const expected = `a PAJLADA b`

	result, err := Substitute(message, args)
	assertErrorsEqual(t, nil, err)
	assertStringsEqual(t, expected, result)
}

func TestMultipleFilters(t *testing.T) {
	const message = `a $(user.name|toupper|tolower) b $(user.level)`
	const expected = `a pajlada b 100`
	// const message = `a $(user.name|toupper|tolower) b`
	// const expected = `a pajlada b`

	result, err := Substitute(message, args)
	assertErrorsEqual(t, nil, err)
	assertStringsEqual(t, expected, result)
}

func TestThreeCommands(t *testing.T) {
	const message = `a $(user.name|toupper|tolower) b $(user.level) c $(user.level)`
	const expected = `a pajlada b 100 c 100`
	// const message = `a $(user.name|toupper|tolower) b`
	// const expected = `a pajlada b`

	result, err := Substitute(message, args)
	assertErrorsEqual(t, nil, err)
	assertStringsEqual(t, expected, result)
}

func TestNonExistantFilter(t *testing.T) {
	const message = `a $(user.name|nonexistantfilter) b`

	result, err := Substitute(message, args)
	assertErrorsEqual(t, ErrNonExistantFilter, err)
	assertStringsEqual(t, message, result)
}

func TestNonExistantArgument(t *testing.T) {
	const message = `a $(nonexistantargument.name|nonexistantfilter) b`

	result, err := Substitute(message, args)
	assertErrorsEqual(t, ErrNonExistantArgument, err)
	assertStringsEqual(t, message, result)
}

func TestNonExistantKey(t *testing.T) {
	const message = `a $(user.nonexistantkey|tolower) b`
	const expected = `a  b`

	result, err := Substitute(message, args)
	assertErrorsEqual(t, nil, err)
	assertStringsEqual(t, expected, result)
}

func TestNoArgumentsProvided(t *testing.T) {
	const message = `a $(nonexistantargument.name|nonexistantfilter) b`

	result, err := Substitute(message, nil)
	assertErrorsEqual(t, ErrNoArgumentsProvided, err)
	assertStringsEqual(t, message, result)
}
