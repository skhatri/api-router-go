package router

import (
	"encoding/json"
	"github.com/skhatri/api-router-go/router/model"
	"github.com/skhatri/api-router-go/router/settings"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type httpRouter struct {
	router  *HttpRouterConfiguration
	options *HttpRouterOptions
}
type RequestSummary struct {
	Method    string `json:"method"`
	Uri       string `json:"uri"`
	TimeTaken int    `json:"time_taken"`
	Unit      string `json:"unit"`
	Status    int    `json:"status_code"`
}

//ServeHTTP http interface method
func (hs *httpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestUri := r.URL.Path
	handlerFunc, pathParams := hs.router.delegate.getHandler(r.Method, requestUri)
	if handlerFunc == nil {
		notFound(&WebRequest{Uri: r.RequestURI})
		return
	}
	var requestData []byte = nil
	if r.Body != nil {
		defer r.Body.Close()
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
		PathParams:  pathParams,
	}

	status := render(w, handlerFunc(webRequest))
	if hs.options.LogRequest {
		timeTaken := time.Since(start).Milliseconds()
		hs.options.LogFunction(RequestSummary{
			Method:    r.Method,
			Uri:       r.RequestURI,
			TimeTaken: int(timeTaken),
			Unit:      "ms",
			Status:    status,
		})
	}
}

func render(w http.ResponseWriter, container *model.Container) int {
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

	encoder := json.NewEncoder(w)
	if container.IsDecorated() {
		encoder.Encode(container)
	} else if status <= 400 {
		encoder.Encode(container.Data)
	} else {
		encoder.Encode(container.Errors)
	}
	return status
}

func Bind(httpRouter *httpRouter, hostPort string) {
	http.ListenAndServe(hostPort, httpRouter)
}
