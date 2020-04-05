[![Build](https://travis-ci.com/skhatri/api-router-go.svg?branch=master)](https://travis-ci.com/github/skhatri/api-router-go)
[![Code Coverage](https://img.shields.io/codecov/c/github/skhatri/api-router-go/master.svg)](https://codecov.io/github/skhatri/api-router-go?branch=master)


### Api-Router-Go

Routing DSL to possibly help in building go API quickly for projects.

#### Running Example App
```
go run router-example.go
```

#### Usage

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
