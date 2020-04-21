package main

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/functions"
	"github.com/skhatri/api-router-go/router/settings"
	"log"
	"net/http"
)

func main() {

	mux := router.NewHttpRouterBuilder().
		WithOptions(router.HttpRouterOptions{
			LogRequest: false,
		}).
		Configure(func(configurer router.ApiConfigurer) {
			_settings := settings.GetSettings()
			configurer.

				Get("/echo", functions.EchoFunc).
				GetIf(true).Register("/status", functions.StatusFunc).
				GetIf(false).Register("/status2", functions.StatusFunc).

				//Style 2
				GetIf(_settings.IsToggleOn("add-more")).
					Add("/status3", functions.StatusFunc).
					Add("/status4", functions.StatusFunc).
				Done().

				Static("test", "test")

		}).Build()
	var address = "0.0.0.0:6100"
	log.Printf("Listening on %s\n", address)
	http.ListenAndServe(address, mux)
}
