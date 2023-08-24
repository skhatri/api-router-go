package starter

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/skhatri/api-router-go/router"
	"github.com/skhatri/api-router-go/router/settings"
)

func parseArguments(args []string, defaultPort int) string {
	app := flag.NewFlagSet("serve", flag.ExitOnError)
	var port = defaultPort
	httpPortFromEnv := os.Getenv("HTTP_PORT")
	if httpPortFromEnv != "" {
		port, _ = strconv.Atoi(httpPortFromEnv)
	}
	var address = "0.0.0.0"
	app.StringVar(&address, "host", address, "Host Interface to listen on")
	app.IntVar(&port, "port", port, "Web port to bind to")
	app.Parse(args[1:])
	if app.Parsed() {
		if port < 1024 || port > 65535 {
			panic("invalid port passed. please provide one between 1024 and 65535")
		}
	}

	return fmt.Sprintf("%s:%d", address, port)
}

//StartApp quick starter
func StartApp(appArgs []string, defaultPort int, configurationHook func(configurer router.ApiConfigurer)) {
	StartAppWithOptions(appArgs, defaultPort, configurationHook, nil)
}

//StartAppWithOptions quick starter that takes a logging function
func StartAppWithOptions(appArgs []string, defaultPort int, configurationHook func(router.ApiConfigurer), logFn func(summary router.RequestSummary)) {
	var args []string

	if len(appArgs) < 2 {
		args = []string{
			"serve",
		}
	} else {
		args = appArgs[1:]
	}
	var command = args[0]

	switch command {
	case "serve":
		address := parseArguments(args, defaultPort)
		var logFunc func(summary router.RequestSummary)
		if logFn == nil {
			logFunc = func(values router.RequestSummary) {
				log.Println(values)
			}
		} else {
			logFunc = logFn
		}
		mux := router.NewHttpRouterBuilder().
			WithOptions(router.HttpRouterOptions{
				LogRequest:  true,
				LogFunction: logFunc,
			}).Configure(func(configurer router.ApiConfigurer) {
			configurationHook(configurer)
		}).Build()
		log.Printf("Listening on %s\n", address)
		http.ListenAndServe(address, mux)
	default:
		log.Printf("command %s is not supported\n", command)
	}

}

//RunApp - loads config from ROUTE_SETTINGS file or router.json
func RunApp(configFn func(router.ApiConfigurer)) {
	_settings := settings.GetSettings()
	mux := router.NewHttpRouterBuilder().
		WithOptions(router.HttpRouterOptions{
			LogRequest: _settings.GetBoolOption("log-request", true),
		}).
		Configure(configFn).Build()

	port := _settings.Transport().Port
	var address = fmt.Sprintf("0.0.0.0:%d", port)

	tls := _settings.Transport().Tls
	if tls != nil && tls.Enabled != nil && *tls.Enabled {
		log.Printf("Listening with TLS on %s\n", address)
		http.ListenAndServeTLS(address, tls.PublicKey, tls.PrivateKey, mux)
	} else {
		log.Printf("Listening on %s\n", address)
		http.ListenAndServe(address, mux)
	}
}
