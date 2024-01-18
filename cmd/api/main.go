package main

import (
	"fmt"
	"net/http"

	routey "github.com/joseph-beck/routey/pkg/router"
	swaggy "github.com/joseph-beck/routey/pkg/swagger"
	swaggerFiles "github.com/swaggo/files"
)

func health() routey.HandlerFunc {
	return func(c *routey.Context) {
		c.Render(http.StatusOK, "health")
	}
}

func ping() routey.HandlerFunc {
	type Thing struct {
		Name  string `json:"name"`
		Other string `json:"other"`
	}

	f := func(c *routey.Context) {
		var t Thing
		c.ShouldBindJSON(&t)
		fmt.Println(t)
	}

	return func(c *routey.Context) {
		f(c)

		var t Thing
		c.ShouldBindJSON(&t)
		fmt.Println(t)
		p := c.Protocol()

		c.Render(http.StatusOK, "pong "+p)
	}
}

func hello() routey.HandlerFunc {
	return func(c *routey.Context) {
		b, err := c.Body()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b) + "\nbody")

		c.Status(http.StatusOK)
	}
}

func middleware() routey.MiddlewareFunc {
	return func(c *routey.Context) {
		fmt.Println("middleware")
	}
}

func shutdown() routey.ShutdownFunc {
	return func() {
		fmt.Println("shutting down")
	}
}

// @title routey
// @version 0.0.1
// @description routey docs teser.

// @license.name MIT

// @host localhost:3000
// @BasePath /api/v1
func main() {
	c := routey.Config{
		Port:  ":3000",
		Debug: true,
		CORS:  false,
	}
	r := routey.New(c)
	r.Use(middleware())
	url := swaggy.URL("http://localhost:3000/docs/doc.json")
	r.Docs("/docs/*", swaggy.WrapHandler(swaggerFiles.Handler, url))
	r.Add(routey.Get, "/api/health", "", health(), nil)
	r.Add(routey.Get, "/api/ping", "", ping(), nil)
	r.Add(routey.Get, "/api/hello", "", hello(), nil)
	go r.Shutdown(shutdown())
	r.Run()
}
