package main

import (
	"fmt"

	routey "github.com/joseph-beck/routey/pkg/router"
	"github.com/joseph-beck/routey/pkg/status"
)

func reverseString(input string) string {
	runes := []rune(input)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func main() {
	r := routey.New()
	go r.Shutdown()
	// Test route that returns some http.

	r.Route(routey.Route{
		Path:   "/",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			b := "<h1>Healthy</h1>"
			c.Render(status.OK, b)
		},
		DecoratorFunc: func(f routey.HandlerFunc) routey.HandlerFunc {
			return func(c *routey.Context) {
				fmt.Println("hello world!")
				f(c)
			}
		},
	})
	// Test error route
	r.Route(routey.Route{
		Path:   "/400",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			c.Status(status.BadRequest)
		},
		DecoratorFunc: nil,
	})
	// Test route that panics
	r.Route(routey.Route{
		Path:   "/panic",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			panic("this is a panic")
		},
		DecoratorFunc: nil,
	})
	// Test route that has parameters
	r.Route(routey.Route{
		Path:   "/echo",
		Params: "/:string",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			p, err := c.Param("string")
			if err != nil {
				c.Status(status.BadRequest)
				return
			}
			c.Render(status.OK, p)
		},
		DecoratorFunc: nil,
	})
	// Test param int
	r.Route(routey.Route{
		Path:   "/add",
		Params: "/:one/:two",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			o, err := c.ParamInt("one")
			if err != nil {
				c.Status(status.BadRequest)
				return
			}
			t, err := c.ParamInt("two")
			if err != nil {
				c.Status(status.BadRequest)
				return
			}
			a := o + t

			c.Render(status.OK, fmt.Sprintf("%d", a))
		},
		DecoratorFunc: nil,
	})
	// Test JSON
	r.Route(routey.Route{
		Path:   "/json",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			type json struct {
				Name string `json:"name"`
				Body string `json:"body"`
			}

			j := json{
				Name: "hello",
				Body: "these are a lot of words.",
			}

			c.JSON(status.OK, j)
		},
		DecoratorFunc: nil,
	})
	// Test Query
	r.Route(routey.Route{
		Path:   "/query",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			w, err := c.Query("word")
			if err != nil {
				c.Status(status.BadRequest)
				return
			}

			c.Render(status.OK, w)
		},
		DecoratorFunc: nil,
	})
	// Test JSON binding
	r.Route(routey.Route{
		Path:   "/bind",
		Params: "",
		Method: routey.Post,
		HandlerFunc: func(c *routey.Context) {
			type obj struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			}

			var o obj
			err := c.ShouldBindJSON(&o)
			if err != nil {
				c.Status(status.BadRequest)
			}

			fmt.Printf("%s name \n%s email\n", o.Name, o.Email)

			c.JSON(status.OK, o)
		},
		DecoratorFunc: nil,
	})
	// Test params and query
	r.Add(
		routey.Get,
		"/reverse",
		"/:text",
		func(c *routey.Context) {
			t, err := c.Param("text")
			if err != nil {
				c.Status(status.BadRequest)
				return
			}
			r, err := c.QueryInt("repeat")
			if err != nil {
				c.Render(status.OK, reverseString(t))
				return
			}

			b := ""
			for i := 0; i < r; i++ {
				b += reverseString(t)
			}
			c.Render(status.OK, b)
		},
		nil,
	)
	// Test HTML files
	r.LoadHTMLGlob("web/*")
	r.Add(
		routey.Get,
		"/html",
		"",
		func(c *routey.Context) {
			c.HTML(
				status.OK,
				"index.html",
				routey.M{
					"title": "Main website",
				},
			)
		},
		nil,
	)
	r.Add(
		routey.Get,
		"/hello",
		"",
		func(c *routey.Context) {
			c.HTML(
				status.OK,
				"hello.html",
				nil,
			)
		},
		nil,
	)
	// Test redirect
	r.Add(
		routey.Get,
		"/redirect",
		"",
		func(c *routey.Context) {
			c.Redirect(status.MovedPermanently, "https://www.google.com")
		},
		nil,
	)

	r.Run()
}
