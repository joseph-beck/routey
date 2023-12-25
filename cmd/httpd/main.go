package main

import (
	"fmt"
	"net/http"

	routey "github.com/joseph-beck/routey/pkg/router"
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

		c.Render(http.StatusOK, "pong")
	}
}

func main() {
	c := routey.Config{
		Port:  ":3000",
		Debug: true,
		CORS:  false,
	}
	r := routey.New(c)
	r.Add(routey.Get, "/api/health", "", health(), nil)
	r.Add(routey.Get, "/api/ping", "", ping(), nil)
	r.Run()
}
