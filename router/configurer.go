package router

import (
	"fmt"
	"github.com/skhatri/api-router-go/router/settings"
	"log"
	"net/http"
)

type HttpRouterConfiguration struct {
	delegate *httpRouterDelegate
}

//HttpRouterOptions Router Options
type HttpRouterOptions struct {
	LogRequest  bool
	LogFunction func(...interface{})
}

type HttpRouteBuilder interface {
	WithOptions(options HttpRouterOptions) HttpRouteBuilder
	Configure(configurerFn func(httpDelegate ApiConfigurer)) HttpRouteBuilder
	SettingsFrom(settingsFile *string) HttpRouteBuilder
	Build() http.Handler
}

type _HttpRouterBuilder struct {
	options         *HttpRouterOptions
	configuration   *HttpRouterConfiguration
	configurationFn func(configurer ApiConfigurer)
	settings        *string
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

func (hrb *_HttpRouterBuilder) SettingsFrom(settingsFile *string) HttpRouteBuilder {
	hrb.settings = settingsFile
	return hrb
}

func (hrb *_HttpRouterBuilder) Build() http.Handler {
	if hrb.options == nil {
		var defaultOptions = HttpRouterOptions{
			LogRequest:  true,
			LogFunction: DefaultLogger,
		}
		hrb.options = &defaultOptions
	}

	hrb.processExternalConfiguration()

	var apiConfigurer ApiConfigurer = hrb.configuration.delegate

	if hrb.configurationFn != nil {
		hrb.configurationFn(apiConfigurer)
	}
	return &httpRouter{
		options: hrb.options,
		router:  hrb.configuration,
	}
}

func (hrb *_HttpRouterBuilder) processExternalConfiguration() {
	err := settings.ApplySettings(hrb.settings)
	if err != nil {
		panic(fmt.Sprintf("error processing route settings %s", err.Error()))
	}
}

func NewHttpRouterBuilder() HttpRouteBuilder {
	return &_HttpRouterBuilder{
		configuration: defaultRouterConfiguration(),
	}
}

func defaultRouterConfiguration() *HttpRouterConfiguration {
	routerRef := &httpRouterDelegate{
		mapping: make(map[string]HandlerFunc),
	}
	return &HttpRouterConfiguration{delegate: routerRef}
}

var DefaultLogger = func(s ...interface{}) {
	log.Println(s)
}
