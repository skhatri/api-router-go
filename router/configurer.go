package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/skhatri/api-router-go/router/settings"
)

type HttpRouterConfiguration struct {
	delegate *httpRouterDelegate
}

//HttpRouterOptions Router Options
type HttpRouterOptions struct {
	LogRequest  bool
	LogFunction func(RequestSummary)
}

type HttpRouteBuilder interface {
	WithOptions(options HttpRouterOptions) HttpRouteBuilder
	Configure(configurerFn func(httpDelegate ApiConfigurer)) HttpRouteBuilder
	Build() *http.ServeMux
}

type _HttpRouterBuilder struct {
	options         *HttpRouterOptions
	configuration   *HttpRouterConfiguration
	configurationFn func(configurer ApiConfigurer)
}

func (hrb *_HttpRouterBuilder) WithOptions(options HttpRouterOptions) HttpRouteBuilder {
	var defaultOptions = options
	if defaultOptions.LogFunction == nil {
		defaultOptions.LogFunction = DefaultLogger
	}
	hrb.options = &defaultOptions
	return hrb
}

func (hrb *_HttpRouterBuilder) Configure(configurerFn func(httpDelegate ApiConfigurer)) HttpRouteBuilder {
	hrb.configurationFn = configurerFn
	return hrb
}

func (hrb *_HttpRouterBuilder) Build() *http.ServeMux {

	if hrb.options == nil {
		var defaultOptions = HttpRouterOptions{
			LogRequest:  true,
			LogFunction: DefaultLogger,
		}
		hrb.options = &defaultOptions
	}

	var apiConfigurer ApiConfigurer = hrb.configuration.delegate

	if hrb.configurationFn != nil {
		hrb.configurationFn(apiConfigurer)
	}
	serveMux := http.NewServeMux()
	addStaticMapping(serveMux, hrb.configuration.delegate.staticMapping)
	addStaticMapping(serveMux, settings.GetSettings().StaticMappings())

	serveMux.Handle("/", hrb.newRouter())
	return serveMux
}

func addStaticMapping(serveMux *http.ServeMux, mapping map[string]string) {
	for pathUri, folder := range mapping {
		fs := http.FileServer(http.Dir(fmt.Sprintf("./%s", folder)))
		prefix := fmt.Sprintf("/%s/", pathUri)
		serveMux.Handle(prefix, http.StripPrefix(prefix, fs))
	}
}

func (hrb *_HttpRouterBuilder) newRouter() *httpRouter {
	return &httpRouter{
		options: hrb.options,
		router:  hrb.configuration,
	}
}

func NewHttpRouterBuilder() HttpRouteBuilder {
	return &_HttpRouterBuilder{
		configuration: defaultRouterConfiguration(),
	}
}

func defaultRouterConfiguration() *HttpRouterConfiguration {
	routerRef := &httpRouterDelegate{
		dynamicStore:  pathMatchingUriStore(),
		staticStore:   simpleUriStore(),
		staticMapping: make(map[string]string),
	}
	return &HttpRouterConfiguration{delegate: routerRef}
}

var DefaultLogger = func(s RequestSummary) {
	log.Println(s.Status, s.Method, s.Uri, s.TimeTaken, s.Unit)
}
