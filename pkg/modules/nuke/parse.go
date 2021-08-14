package nuke

import (
	"regexp"
	"strings"
	"time"
)

type NukeParameterParser struct {
}

func (p *NukeParameterParser) ParseNukeParameters(parts []string) (*NukeParameters, error) {
	if len(parts) < 3 {
		return nil, ErrUsage
	}

	phrase := strings.Join(parts[1:len(parts)-2], " ")
	var regexPhrase *regexp.Regexp
	if strings.HasPrefix(phrase, "/") && strings.HasSuffix(phrase, "/") {
		// parse as regex
		asd := phrase[1 : len(phrase)-1]
		regex, err := regexp.Compile(asd)
		if err == nil {
			regexPhrase = regex
		}
	}

	scrollbackLength, err := time.ParseDuration(parts[len(parts)-2])
	if err != nil {
		return nil, ErrUsage
	}
	if scrollbackLength < 0 {
		return nil, ErrUsage
	}
	timeoutDuration, err := time.ParseDuration(parts[len(parts)-1])
	if err != nil {
		return nil, ErrUsage
	}
	if timeoutDuration < 0 {
		return nil, ErrUsage
	}

	return &NukeParameters{
		Phrase:           phrase,
		RegexPhrase:      regexPhrase,
		ScrollbackLength: scrollbackLength,
		TimeoutDuration:  timeoutDuration,
	}, nil
}
