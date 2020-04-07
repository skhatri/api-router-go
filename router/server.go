package router

import (
	"encoding/json"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/api-router-go/router/settings"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpRouter struct {
	router  *HttpRouterConfiguration
	options *HttpRouterOptions
}

//ServeHTTP http interface method
func (hs *httpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUri := r.URL.Path
	handlerFunc := hs.router.delegate.getHandler(r.Method, requestUri)
	if handlerFunc == nil {
		notFound(&WebRequest{Uri: r.RequestURI})
		return
	}
	if hs.options.LogRequest {
		hs.options.LogFunction("recv", r.RequestURI)
	}
	var requestData []byte = nil
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			render(w, badRequest())
			return
		}
		requestData = body
	}
	var headers = make(map[string]string)
	for name, value := range r.Header {
		headers[name] = strings.Join(value, ";")
	}
	var params = make(map[string]string)

	for name, value := range r.URL.Query() {
		params[name] = strings.Join(value, ";")
	}

	webRequest := &WebRequest{
		Uri:         requestUri,
		Body:        requestData,
		Headers:     headers,
		QueryParams: params,
		QueryString: r.URL.RawQuery,
	}
	render(w, handlerFunc(webRequest))
}

func render(w http.ResponseWriter, container *model.Container) {
	status := container.GetStatus()
	if status == 0 {
		status = 200
	}
	for k, v := range settings.GetSettings().ResponseHeaders() {
		w.Header().Add(k, v)
	}
	for k, v := range container.GetHeaders() {
		w.Header().Add(k, v)
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(container)
}

func Bind(httpRouter *httpRouter, hostPort string) {
	http.ListenAndServe(hostPort, httpRouter)
}
