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