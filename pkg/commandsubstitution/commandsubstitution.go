package commandsubstitution

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Filter interface {
	Run(string) string
}

type filterToLower struct {
}

func (filterToLower) Run(in string) string {
	return strings.ToLower(in)
}

type filterToUpper struct {
}

func (filterToUpper) Run(in string) string {
	return strings.ToUpper(in)
}

const (
	baseRegex = `\$\((?P<command>[a-zA-Z\.]+)((?P<namedargument>:)|(?P<argument>;[0-9]+))?((?:\|\w+)*)\)`
	// baseRegex     = `e`
	urlfetchRegex = `$\(urlfetch (?P<url>[a-zA-Z0-9~:/?#\[\]@!$%\'()*+,;=\.]+)\)`
)

var (
	filters = map[string]Filter{}

	ErrNonExistantFilter = errors.New("attempted to use a filter that is not registered")

	ErrNonExistantArgument = errors.New("attempted to use an argument that wasn't provided")

	ErrNoArgumentsProvided = errors.New("no arguments were provided, so no substitutions can be made")
)

func init() {
	filters["tolower"] = &filterToLower{}
	filters["toupper"] = &filterToUpper{}
}

type Substitution interface {
	GetKey(key string) (value string)
}

type User struct {
	name  string
	level int
}

// maybe return error as second value
func (u User) GetKey(key string) string {
	switch key {
	case "name":
		return u.name
	case "level":
		return strconv.FormatInt(int64(u.level), 10)
	}

	return ""
}

// Substitute attempts to substitute all command substitutions (i.e. $(user.name)) in `message` with readable values
func Substitute(message string, arguments map[string]Substitution) (string, error) {
	if arguments == nil {
		return message, ErrNoArgumentsProvided
	}

	fullResult := message

	r := regexp.MustCompile(baseRegex)
	indices := r.FindAllStringSubmatchIndex(message, -1)
	offset := 0
	for _, matchIndices := range indices {
		command := message[matchIndices[2]:matchIndices[3]]
		commandParts := strings.Split(command, ".")
		substitution, ok := arguments[commandParts[0]]
		if !ok {
			return message, ErrNonExistantArgument
		}
		result := substitution.GetKey(commandParts[1])
		if matchIndices[10] != -1 {
			subFilters := strings.Split(message[matchIndices[10]:matchIndices[11]], "|")[1:]
			if len(subFilters) > 0 {
				for _, filterString := range subFilters {
					filter, ok := filters[filterString]
					if !ok {
						return message, ErrNonExistantFilter
					}

					result = filter.Run(result)
				}
			}
		}

		fullResult = fullResult[:matchIndices[0]-offset] + result + fullResult[matchIndices[1]-offset:]
		offset = len(message) - len(fullResult)
	}

	return fullResult, nil
}
