package model

import (
	"regexp"
)

type UriMatcher interface {
	Matches(uri string) bool
}

type RegexMatcher struct {
	RegexPattern *regexp.Regexp
}

func (re *RegexMatcher) Matches(uri string) bool {
	str := re.RegexPattern.FindStringSubmatch(uri)
	return len(str) > 0
}
