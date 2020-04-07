package starter

import (
	"flag"
	"fmt"
	"github.com/skhatri/api-router-go/router"
	"log"
	"net/http"
	"os"
	"strconv"
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

//Application Quickstart
func StartApp(appArgs []string, defaultPort int, configurationHook func(router.ApiConfigurer)) {
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
		r := router.NewHttpRouterBuilder().
			WithOptions(router.HttpRouterOptions{
				LogRequest: true,
				LogFunction: func(values ...interface{}) {
					log.Println(values)
				},
			}).Configure(func(configurer router.ApiConfigurer) {
			configurationHook(configurer)
		}).Build()
		log.Printf("Listening on %s\n", address)
		http.ListenAndServe(address, r)
	default:
		log.Printf("command %s is not supported\n", command)
	}

}
