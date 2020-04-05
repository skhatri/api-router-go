package router

import "fmt"

type HttpRouterConfiguration struct {
	delegate *httpRouterDelegate
}

//HttpRouterOptions Router Options
type HttpRouterOptions struct {
	LogRequest  bool
	LogFunction func(...interface{})
}

type HttpRouterBuilder struct {
	options       *HttpRouterOptions
	configuration *HttpRouterConfiguration
}

func (hrb *HttpRouterBuilder) WithOptions(options HttpRouterOptions) *HttpRouterBuilder {
	var defaultOptions = options
	if defaultOptions.LogFunction == nil {
		defaultOptions.LogFunction = DefaultLogger
	}
	hrb.options = &defaultOptions
	return hrb
}

func (hrb *HttpRouterBuilder) Configure(configurerFn func(httpDelegate ApiConfigurer)) *HttpRouterBuilder {
	var apiConfigurer ApiConfigurer = hrb.configuration.delegate
	configurerFn(apiConfigurer)
	return hrb
}

func (hrb *HttpRouterBuilder) Build() *httpRouter {
	if hrb.options == nil {
		var defaultOptions = HttpRouterOptions{
			LogRequest:  true,
			LogFunction: DefaultLogger,
		}
		hrb.options = &defaultOptions
	}
	return &httpRouter{
		options: hrb.options,
		router:  hrb.configuration,
	}
}

func NewHttpRouterBuilder() *HttpRouterBuilder {
	return &HttpRouterBuilder{
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
	fmt.Println(s)
}
