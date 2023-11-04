package router

import (
	"log"
	"net/http"

	errs "github.com/joseph-beck/routey/pkg/error"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func (c *Context) Error(err errs.Error) {
	log.Println(err.String())
}

func (c *Context) Write(s string) (int, error) {
	b, err := c.w.Write([]byte(s))
	return b, err
}

func (c *Context) Header(k string, v string) {
	if v == "" {
		c.w.Header().Del(k)
		return
	}

	if k == "" {
		return
	}

	c.w.Header().Set(k, v)
}

func (c *Context) GetHeader(k string) string {
	v := c.r.Header.Get(k)
	return v
}

func (c *Context) Status(s int) {
	c.w.WriteHeader(s)
}

func (c *Context) Render(s int, b string) {
	c.Status(s)

	_, err := c.Write(b)
	if err != nil {
		c.Error(errs.RenderError)
	}
}
