package functions

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/model"
)

var statusResult = map[string]string{
	"status": "OK",
}

func StatusFunc(_ *router.WebRequest) *model.Container {
	return model.WithDataOnly(statusResult)
}

func EchoFunc(request *router.WebRequest) *model.Container {
	container := model.Response(request)
	for k, v := range request.Headers {
		container.AddHeader(k, v)
	}
	return container
}

func IndexFunc(_ *router.WebRequest) *model.Container {
	return model.Response(map[string]string{
		"message": "It works",
	})
}
