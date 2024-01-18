package main

import (
	"fmt"
	"net/http"

	routey "github.com/joseph-beck/routey/pkg/router"
	swaggy "github.com/joseph-beck/routey/pkg/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// health godoc
// @Summary      API Health
// @Description  Get the Health of the API
// @Tags         health
// @Success      200
// @Router       /api/health [get]
func health() routey.HandlerFunc {
	return func(c *routey.Context) {
		c.Render(http.StatusOK, "health")
	}
}

// ping godoc
// @Summary      API Ping
// @Description  Ping the API and receive a JSON response
// @Tags         ping
// @Produce 	 json
// @Success      200
// @Router       /api/ping [get]
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

// hello godoc
// @Summary      API Hello
// @Description  Say hello to the API
// @Tags         health
// @Success      200
// @Router       /api/hello [get]
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

	}
}

func shutdown() routey.ShutdownFunc {
	return func() {
		fmt.Println("shutting down")
	}
}

// @title routey
// @version 0.0.1
// @description routey docs tester.
// @termsOfService http://swagger.io/terms/

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
	config := swaggy.URL("/docs/swagger.json")
	r.Docs("/docs/*", swaggy.WrapHandler(swaggerFiles.Handler, config))
	r.Add(routey.Get, "/api/health", "", health(), nil)
	r.Add(routey.Get, "/api/ping", "", ping(), nil)
	r.Add(routey.Get, "/api/hello", "", hello(), nil)
	go r.Shutdown(shutdown())
	r.Run()
}
