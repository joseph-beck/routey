# routey

An extremely simple Go HTTP Router made for fun!

## Examples

### Getting Started

```go
package main

import (
    "fmt"
    "net/http"

    routey "github.com/joseph-beck/routey/pkg/router"
)

func main() {
    r := routey.New()
    go r.Shutdown()

    r.Route(routey.Route{
        Path: "/hello",
        Params: "",
        Method: routey.Get,
        HandleFunc: func(c *routey.Context) {
            b := "Hello world!"
            c.Render(http.StatusOK, b)
        }
        DecoratorFunc: nil,
    })

    r.Run()
}

```

### Using parameters

```go
func handler() routey.HandlerFunc {
    return func(c *routey.Context) {
        p, err := c.Param("param")
        if err != nil {
            c.Status(http.BadRequest)
            return
        }

        c.Render(http.StatusOK, p)
    }
}

r.Add(
    "/route",
    "/:param",
    routey.Get,
    handler(),
    nil
)
```

### Using queries

```go
func handler() routey.HandlerFunc {
    return func(c *routey.Context) {
        q, err := c.Query("query")
        if err != nil {
            c.Status(http.BadRequest)
            return
        }

        c.Render(http.StatusOK, q)
    }
}
```

### Creating a HandlerFunc

```go
func handler() routey.HandlerFunc {
    return func(c *routey.Context) {
        b := "I am a handler!"
        c.Render(http.StatusOK, b)
    }
}
```

Declaring a handler function this way allows us to more easily use dependency injection.

### Creating a DecoratorFunc

```go
func decorator(f routey.HandlerFunc) routey.HandlerFunc {
    return func(c *routey.Context) {
        fmt.Println("I am a decorator!")
        f(c)
    }
}
```

Declaring a decorator function this way allows us to decorate decorator functions as well as more easily use dependency injection. They can be used for a variety of things, but commonly used in protecting our end points.
