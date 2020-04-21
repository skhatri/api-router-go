package router

import (
	"fmt"
	"github.com/skhatri/api-router-go/router/model"
	"strings"
)

type WebRequest struct {
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query-params"`
	Body        []byte            `json:"body"`
	Uri         string            `json:"uri"`
	QueryString string            `json:"query"`
}

func (web *WebRequest) GetHeader(name string) string {
	return web.Headers[name]
}

func (web *WebRequest) GetQueryParam(name string) string {
	return web.QueryParams[name]
}

type HandlerFunc func(request *WebRequest) *model.Container

type ApiConfigurer interface {
	Get(string, HandlerFunc) ApiConfigurer
	Post(string, HandlerFunc) ApiConfigurer
	Method(string, string, HandlerFunc) ApiConfigurer
	GetIf(cond bool) *ConditionalMethodBuilder
	PostIf(cond bool) *ConditionalMethodBuilder
	Static(string, string) ApiConfigurer
}

type httpRouterDelegate struct {
	mapping map[string]HandlerFunc
	staticMapping map[string]string
}

func (router *httpRouterDelegate) Get(uri string, handlerFn HandlerFunc) ApiConfigurer {
	router.Method("GET", uri, handlerFn)
	return router
}

func (router *httpRouterDelegate) Static(path string, folder string) ApiConfigurer {
	router.staticMapping[path] = folder
	return router
}

type ConditionalMethodBuilder struct {
	Method string
	Check    bool
	Delegate *httpRouterDelegate
}

func (methodBuilder *ConditionalMethodBuilder) Register(uri string, handlerFunc HandlerFunc) ApiConfigurer {
	if methodBuilder.Check {
		methodBuilder.Delegate.Method(methodBuilder.Method, uri, handlerFunc)
	}
	return methodBuilder.Delegate
}

func (methodBuilder *ConditionalMethodBuilder) Add(uri string, handlerFunc HandlerFunc) *ConditionalMethodBuilder {
	if methodBuilder.Check {
		methodBuilder.Delegate.Method(methodBuilder.Method, uri, handlerFunc)
	}
	return methodBuilder
}

func (methodBuilder *ConditionalMethodBuilder) Done() ApiConfigurer {
	return methodBuilder.Delegate
}

func (router *httpRouterDelegate) GetIf(cond bool) *ConditionalMethodBuilder {
	return &ConditionalMethodBuilder{
		Method: "GET",
		Check:    cond,
		Delegate: router,
	}
}

func (router *httpRouterDelegate) PostIf(cond bool) *ConditionalMethodBuilder {
	return &ConditionalMethodBuilder{
		Method: "POST",
		Check:    cond,
		Delegate: router,
	}
}

func (router *httpRouterDelegate) Post(uri string, handlerFn HandlerFunc) ApiConfigurer {
	router.Method("POST", uri, handlerFn)
	return router
}

func (router *httpRouterDelegate) Method(method string, uri string, handlerFunc HandlerFunc) ApiConfigurer {
	methodKey := strings.ToLower(fmt.Sprintf("%s::%s", method, uri))
	router.mapping[methodKey] = handlerFunc
	return router
}

func (router *httpRouterDelegate) getHandler(method string, uri string) HandlerFunc {
	methodKey := strings.ToLower(fmt.Sprintf("%s::%s", method, uri))
	handlerFunc, ok := router.mapping[methodKey]
	if ok {
		return handlerFunc
	}
	return notFound
}

func notFound(request *WebRequest) *model.Container {
	return model.ErrorResponse(model.MessageItem{
		Code:    "not-found",
		Message: fmt.Sprintf("uri %s not found", request.Uri),
		Details: nil,
	}, 404)
}

func badRequest() *model.Container {
	return model.ErrorResponse(model.MessageItem{
		Code:    "Bad Request",
		Message: "Bad Request",
		Details: nil,
	}, 400)
}

func internalError() *model.Container {
	return model.ErrorResponse(model.MessageItem{
		Code:    "Internal Error",
		Message: "Internal Error",
		Details: nil,
	}, 500)
}
