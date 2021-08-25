package router

import (
	"fmt"
	"github.com/skhatri/api-router-go/router/model"
	"strings"
)

type WebRequest struct {
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query-params"`
	PathParams  map[string]string `json:"path-params"`
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

func (web *WebRequest) GetPathParam(name string) string {
	return web.PathParams[name]
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
	staticMapping map[string]string
	staticStore   HandlerRegistry
	dynamicStore  HandlerRegistry
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
	Method   string
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
		Method:   "GET",
		Check:    cond,
		Delegate: router,
	}
}

func (router *httpRouterDelegate) PostIf(cond bool) *ConditionalMethodBuilder {
	return &ConditionalMethodBuilder{
		Method:   "POST",
		Check:    cond,
		Delegate: router,
	}
}

func (router *httpRouterDelegate) Post(uri string, handlerFn HandlerFunc) ApiConfigurer {
	router.Method("POST", uri, handlerFn)
	return router
}

func (router *httpRouterDelegate) Method(method string, uri string, handlerFunc HandlerFunc) ApiConfigurer {
	if strings.Contains(uri, ":") {
		router.dynamicStore.Add(method, uri, handlerFunc)
		router.dynamicStore.Add("OPTIONS", uri, handlerFunc)
	} else {
		router.staticStore.Add(method, uri, handlerFunc)
		router.staticStore.Add("OPTIONS", uri, handlerFunc)
	}
	return router
}

func (router *httpRouterDelegate) getHandler(method string, uri string) (HandlerFunc, bool, map[string]string) {
	handlerFunc, _ := router.staticStore.Lookup(method, uri)
	if handlerFunc != nil {
		return handlerFunc, true, nil
	}
	handlerFunc, params := router.dynamicStore.Lookup(method, uri)
	if handlerFunc != nil {
		return handlerFunc, true, params
	}
	return notFound, false, nil
}

func ok(request *WebRequest) *model.Container {
	return model.ResponseWithStatusCode(make(map[string]interface{}, 0), 200)
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
