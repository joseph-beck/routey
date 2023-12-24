# routey

An extremely simple Go HTTP Router made for fun! This was heavily inspired by gin.

## Installation

```sh
$ go get -u github.com/joseph-beck/routey
```

```go
import (
    routey github.com/joseph-beck/routey/pkg/router
)
```

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
    nil,
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

### Services

routey also supports the use of Services, which are structs with methods that are your endpoints, every Service must implement an `Add() []Route` that can be registered to the App with the `Register()` method.

```go
type HelloService struct{}

func (s HelloService) Add() []Route {
    return []Route{
        {
            Path:          "/hello",
            Params:        "",
            Method:        Get,
            HandlerFunc:   s.Get(),
            DecoratorFunc: nil,
        },
    }
}

func (s *HelloService) Get() HandlerFunc {
    return func(c *Context) {
        c.Status(http.StatusOK)
    }
}
```

When creating instances of your Service you can give dependencies to the struct for example a database connection, etc.

### Rendering HTML

```go
func main() {
    ...
    r.LoadHTMLGlob("web/*")
    ...
}

// /index
func indexHandler() routey.HandlerFunc {
    return func(c *routey.Context) {
        c.HTML(
            http.StatusOK,
            "index.html",
            nil,
        )
    }
}

// /name/:name
func helloHandler() routey.HandlerFunc {
    return func(c *routey.Context) {
        name, err := c.Param("name")
        if err != nil {
            c.Status(http.StatusBadRequest)
            return
        }

        c.HTML(
            http.StatusOK,
            "hello.html",
            routey.M{
                "name": name,
            },
        )
    }
}
```

We can render HTML pages by loading a glob, or group of HTML files in a directory, or just specific files. The endpoints are added to the App just like any other endpoint. We can link HTML pages together by creating an endpoint for another HTML file and using a href to that, for example:

***index.html***

```html
<!DOCTYPE html>
<html>
<body>

<h1>Routey</h1>
<a href="/hello">hello<a/>

</body>
</html>
```

***hello.html***

```html
<!DOCTYPE html>
<html>
<body>

<p>Hello {{ .name }}<p>
<a href="/index">home<a/>

</body>
</html>

```
