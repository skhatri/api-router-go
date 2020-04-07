package main

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/functions"
	"log"
	"net/http"
)

func main() {
	httpRouter := router.NewHttpRouterBuilder().
		WithOptions(router.HttpRouterOptions{
			LogRequest: false,
		}).
		Configure(func(configurer router.ApiConfigurer) {
			configurer.
				Get("/echo", functions.EchoFunc).
				GetIf(true).Register("/status", functions.StatusFunc).
				GetIf(false).Register("/status2", functions.StatusFunc)
		}).Build()
	var address = "0.0.0.0:6100"
	log.Printf("Listening on %s\n", address)
	http.ListenAndServe(address, httpRouter)
}
