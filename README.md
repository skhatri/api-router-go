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
    import (
        "fmt"
        "github.com/skhatri/api-router-go/router"
        "github.com/skhatri/api-router-go/router/functions"
        "github.com/skhatri/api-router-go/starter"
        "net/http"
    )

    func main() {
      starter.StartApp(os.Args, 6200, func(cfg router.ApiConfigurer) {
        configurer.Get("/echo", functions.EchoFunc)
      })
    }   
```
#### Start your own server 

```go
    import (
        "fmt"
        "github.com/skhatri/api-router-go/router"
        "github.com/skhatri/api-router-go/router/functions"
        "net/http"
    )
    
    func main() {
        httpRouter := router.NewHttpRouterBuilder().
            WithOptions(router.HttpRouterOptions{
                LogRequest: false,
            }).Configure(func(configurer router.ApiConfigurer) {
            configurer.Get("/echo", functions.EchoFunc)
        }).Build()
        var address = "0.0.0.0:6100"
        fmt.Printf("Listening on %s\n", address)
        http.ListenAndServe(address, httpRouter)
    }
```
