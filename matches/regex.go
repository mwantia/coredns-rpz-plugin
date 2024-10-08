package matches

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/coredns/coredns/request"
)

type RegexData struct {
	Entries []struct {
		Pattern string
		Regex   regexp.Regexp
	}
}

func ProcessRegexData(value json.RawMessage) (interface{}, error) {
	var patterns []string
	if err := json.Unmarshal(value, &patterns); err != nil {
		return nil, err
	}

	data := RegexData{}

	for _, pattern := range patterns {
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}

		data.Entries = append(data.Entries, struct {
			Pattern string
			Regex   regexp.Regexp
		}{
			Pattern: pattern,
			Regex:   *regex,
		})
	}

	return data, nil
}

func MatchRegex(state request.Request, ctx context.Context, data RegexData) (*MatchResult, error) {
	qname := state.Name()

	for _, entry := range data.Entries {
		if entry.Pattern == qname {
			return &MatchResult{
				Handled: true,
				Data:    entry.Pattern,
			}, nil
		}

		if entry.Regex.MatchString(qname) {
			return &MatchResult{
				Handled: true,
				Data:    entry.Pattern,
			}, nil
		}
	}

	return nil, nil
}
