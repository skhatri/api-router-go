[![Build](https://travis-ci.com/skhatri/api-router-go.svg?branch=master)](https://travis-ci.com/github/skhatri/api-router-go)
[![Code Coverage](https://img.shields.io/codecov/c/github/skhatri/api-router-go/master.svg)](https://codecov.io/github/skhatri/api-router-go?branch=master)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3825/badge)](https://bestpractices.coreinfrastructure.org/projects/3825)
[![Maintainability](https://api.codeclimate.com/v1/badges/6238e287a522d53ea62c/maintainability)](https://codeclimate.com/github/skhatri/api-router-go/maintainability)

### Api-Router-Go

Routing DSL to possibly help in building go API quickly for projects.

#### Running Example App

```
go run router-example.go
```

#### Quickstart

```go
package main

import (
  "github.com/skhatri/api-router-go/router"
  "github.com/skhatri/api-router-go/router/functions"
  "github.com/skhatri/api-router-go/starter"
  "os"
)

func main() {
  starter.StartApp(os.Args, 6200, func (cfg router.ApiConfigurer) {
    cfg.Get("/echo", functions.EchoFunc)
  })
}

```

#### Start your own server

```go
import (
  "github.com/skhatri/api-router-go/router"
  "github.com/skhatri/api-router-go/starter"
  "github.com/skhatri/api-router-go/router/functions"
)

func main() {
  starter.RunApp(func(configurer router.ApiConfigurer) {
		configurer.Get("/echo", functions.EchoFunc)
	})
}
```

#### Serving Static Paths

Static paths can be provided by ApiConfigurer DSL like this

```
    configurer.
        Static("test", "test").
```

It can also be provided in router.json. Here is the relevant configuration in router.json

```
  "static": {
    "source": ""
  },
```

The above configuration will create static file server to serve all content in the current directory under the path
/source

#### Path Variables Support

Path variables can be specified with a prefix ```:``` in the configured uri.

```
    configurer.
        Get("/greetings/:id", functions.EchoFunc)
``` 

They can be extracted from injected WebRequest using the following:

```
    func myHandler(web *router.WebRequest) *router.Container {
        id := web.GetPathParam("id")
        ...
    }
```

#### Other Configurations

A "router.json" can be provided either in the same folder as the binary or via ROUTE_SETTINGS environment variable to
add default response headers and variables that are to be exposed via RouteSettings.

```
{
  "response-headers": {
    "x-served-by": "Api-Router-Go",
    "access-control-allow-origin": "http://localhost:5000",
    "access-control-allow-methods": "GET, POST, OPTIONS",
    "access-control-allow-headers": "X-Auth-Token, Content-Type, X-Client-Id",
    "access-control-allow-credentials": "false",
    "access-control-max-age": "7200"
  },
  "variables": {
    "service_a_url": "http://localhost:7999"
  },
  "toggles": {
    "read-mode": true,
    "write-mode": false
  },
  "app": {
    
  }
}
```

Variables and Toggles can be read like this

```
settings.GetSettings().Variable("service_a_url");

settings.GetSettings().IsToggleOn("read-mode");
settings.GetSettings().IsToggleOff("write-mode");
settings.GetSettings().IsOn("read-mode");

```

#### Returning Responses

Container type can be returned from each function which decorates with nodes like errors, warning, metadata and data.

```
func StatusFunc(_ *router.WebRequest) *model.Container {
	return model.Response(statusResult)
}
```

The above snippet will render

```
{"data":{"status": "OK"}}
```

If decoration is not preferred, you can create container using ```WithDataOnly```

```
func StatusFunc(_ *router.WebRequest) *model.Container {
	return model.WithDataOnly(statusResult)
}
```

The above snippet will render

```
{"status": "OK"}
```

#### Running TLS

```
openssl genrsa -out private.key 2048
openssl req -new -x509 -sha256 -key private.key -out cert.pem -days 730
```

Add generated key and cert into router.json

```
"transport": {
    "port": 6100,
    "tls": {
      "enabled": true,
      "private-key": "private.key",
      "public-key": "cert.pem"
    }
  }
```


