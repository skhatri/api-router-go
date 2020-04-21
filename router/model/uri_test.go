package model

import (
	"github.com/skhatri/api-router-go/test"
	"regexp"
	"testing"
)

func TestMatchUrlWildcard(t*testing.T) {
	var pattern *regexp.Regexp = regexp.MustCompile("/test")
	var reg = RegexMatcher{
		RegexPattern:pattern,
	}
	tasks := []string {
		"/test",
		"/test/124",
		"/test/123/task",
	};
	for _, task := range tasks {
		test.EqualTo(t, true, reg.Matches(task))
	}

}