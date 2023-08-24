package main

import (
	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/functions"
	"github.com/skhatri/api-router-go/router/settings"
	"github.com/skhatri/api-router-go/starter"
)

func main() {
	_settings := settings.GetSettings()
	addMore := _settings.IsToggleOn("add-more")

	configFn := func(configurer router.ApiConfigurer) {
		configurer.
			Get("/echo", functions.EchoFunc).
			Post("/echo", functions.EchoFunc).
			GetIf(true).Register("/status", functions.StatusFunc).
			GetIf(false).Register("/status2", functions.StatusFunc).

			//Style 2
			GetIf(addMore).
			Add("/status3", functions.StatusFunc).
			Add("/status4", functions.StatusFunc).
			Done().
			Get("/greetings/:id-name", functions.EchoFunc).
			Static("test", "test")

	}
	starter.RunApp(configFn)

}
