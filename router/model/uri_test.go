package model

import (
	"fmt"
	"github.com/skhatri/api-router-go/test"
	"regexp"
	"strings"
	"testing"
)

type PathParam interface {
	ResolveUri(string) map[string]string
}

type PathParamResolver struct {
	regexpPaths []*regexp.Regexp
}

func NewPathParamResolver(registeredPaths []string) PathParamResolver {
	variableName := regexp.MustCompile("(:[a-zA-Z0-9_]+)")
	regexPaths := make([]*regexp.Regexp, 0)
	for _, registeredPath := range registeredPaths {
		newPath := registeredPath
		for _, m := range variableName.FindAllStringSubmatch(registeredPath, -1) {
			newPath = strings.Replace(newPath, m[0], fmt.Sprintf("(?P<%s>[a-zA-Z0-9_]+)", m[0][1:]), -1)
		}
		regexPaths = append(regexPaths, regexp.MustCompile(newPath))
	}
	pathResolver := PathParamResolver{
		regexpPaths: regexPaths,
	}
	return pathResolver
}

func (paramResolver *PathParamResolver) ResolveUri(uri string) map[string]string {
	var paramsMap = make(map[string]string)
	for _, urlPattern := range paramResolver.regexpPaths {
		if urlPattern.Match([]byte(uri)) {
			pathParams := urlPattern.FindAllStringSubmatch(uri, -1)
			var names = urlPattern.SubexpNames()[1:]
			var values = pathParams[0][1:]
			for i := range names {
				paramsMap[names[i]] = values[i]
			}
			break
		}
	}
	return paramsMap
}

func TestMatchUrlWildcard(t *testing.T) {

	pathResolver := NewPathParamResolver([]string{
		"/test/country/:country/id/:id/city/:city",
		"/user/:id",
		"/account/:id/summary",
	})
	paramMap := pathResolver.ResolveUri("/test/country/aus/id/123/city/city_a")
	test.EqualTo(t, "city_a", paramMap["city"])

}
