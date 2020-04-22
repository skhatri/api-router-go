package router

import (
	"fmt"
	"regexp"
	"strings"
)

type HandlerRegistry interface {
	Add(string, string, HandlerFunc)
	Lookup(string, string) (HandlerFunc, map[string]string)
}

type UrlPatternReference struct {
	HandlerKey         string
	PathTemplate       string
	ResolvedExpression string
	PathMatcher        *regexp.Regexp
}
type pathMatchingStore struct {
	rawRef   map[string]HandlerFunc
	patterns map[string][]UrlPatternReference
}

func pathMatchingUriStore() HandlerRegistry {
	init := make(map[string][]UrlPatternReference)
	rawRef := make(map[string]HandlerFunc)
	return &pathMatchingStore{
		patterns: init,
		rawRef:   rawRef,
	}
}

var variableName = regexp.MustCompile("(:[a-zA-Z0-9_]+)")

func pathTemplateToUrlMatcher(registeredPath string) (*regexp.Regexp, string) {
	newPath := registeredPath
	for _, m := range variableName.FindAllStringSubmatch(registeredPath, -1) {
		newPath = strings.Replace(newPath, m[0], fmt.Sprintf("(?P<%s>[a-zA-Z0-9_]+)", m[0][1:]), -1)
	}
	compiled := regexp.MustCompile(newPath)
	return compiled, newPath
}

func (ps *pathMatchingStore) Add(method string, uri string, handler HandlerFunc) {
	methodKey := strings.ToLower(fmt.Sprintf("%s::%s", method, uri))
	ps.rawRef[methodKey] = handler
	uriMatcher, resolvedExpression := pathTemplateToUrlMatcher(uri)
	patternsForMethod := ps.patterns[method]
	if patternsForMethod == nil {
		patternsForMethod = make([]UrlPatternReference, 0)
	}
	patternsForMethod = append(patternsForMethod, UrlPatternReference{
		HandlerKey:         methodKey,
		PathTemplate:       uri,
		PathMatcher:        uriMatcher,
		ResolvedExpression: resolvedExpression,
	})
	ps.patterns[method] = patternsForMethod
}

func (ps *pathMatchingStore) Lookup(method string, uri string) (HandlerFunc, map[string]string) {
	methodRegistry := ps.patterns[method]
	var paramsMap = make(map[string]string)
	var handler HandlerFunc = nil
	if methodRegistry != nil {
		for _, uriReference := range methodRegistry {
			if uriReference.PathMatcher.MatchString(uri) {
				pathParams := uriReference.PathMatcher.FindAllStringSubmatch(uri, -1)
				var names = uriReference.PathMatcher.SubexpNames()[1:]
				var values = pathParams[0][1:]
				for i := range names {
					paramsMap[names[i]] = values[i]
				}
				handler = ps.rawRef[uriReference.HandlerKey]
			}
		}
	}
	return handler, paramsMap
}

type simpleUriMethodStore struct {
	store map[string]HandlerFunc
}

func simpleUriStore() HandlerRegistry {
	init := make(map[string]HandlerFunc)
	return &simpleUriMethodStore{
		store: init,
	}
}

func (ss *simpleUriMethodStore) Add(method string, uri string, handler HandlerFunc) {
	methodKey := strings.ToLower(fmt.Sprintf("%s::%s", method, uri))
	ss.store[methodKey] = handler
}

func (ss *simpleUriMethodStore) Lookup(method string, uri string) (HandlerFunc, map[string]string) {
	methodKey := strings.ToLower(fmt.Sprintf("%s::%s", method, uri))
	fn := ss.store[methodKey]
	return fn, nil
}
