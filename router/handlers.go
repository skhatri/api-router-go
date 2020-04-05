package router

import (
	"fmt"
	"github.com/skhatri/api-router-go/router/model"
	"strings"
)

type WebRequest struct {
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"headers"`
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
	Get(string, HandlerFunc)
	Post(string, HandlerFunc)
	Method(string, string, HandlerFunc)
}

type httpRouterDelegate struct {
	mapping map[string]HandlerFunc
}

func (router *httpRouterDelegate) Get(uri string, handlerFn HandlerFunc) {
	router.Method("GET", uri, handlerFn)
}

func (router *httpRouterDelegate) Post(uri string, handlerFn HandlerFunc) {
	router.Method("POST", uri, handlerFn)
}

func (router *httpRouterDelegate) Method(method string, uri string, handlerFunc HandlerFunc) {
	methodKey := strings.ToLower(fmt.Sprintf("%s::%s", method, uri))
	router.mapping[methodKey] = handlerFunc
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
