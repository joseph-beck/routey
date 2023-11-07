package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	errs "github.com/joseph-beck/routey/pkg/error"
)

type Context struct {
	writer  http.ResponseWriter
	request *http.Request
	params  map[string]string
}

func (c *Context) Error(err errs.Error) {
	log.Println(err.String())
}

func (c *Context) Write(s string) (int, error) {
	b, err := c.writer.Write([]byte(s))
	return b, err
}

func (c *Context) WriteBytes(body []byte) (int, error) {
	b, err := c.writer.Write(body)
	return b, err
}

func (c *Context) Header(k string, v string) {
	if v == "" {
		c.writer.Header().Del(k)
		return
	}

	if k == "" {
		return
	}

	c.writer.Header().Set(k, v)
}

func (c *Context) GetHeader(k string) string {
	v := c.request.Header.Get(k)
	return v
}

func (c *Context) Status(s int) {
	c.writer.WriteHeader(s)
}

func (c *Context) Render(s int, b string) {
	c.Status(s)

	_, err := c.Write(b)
	if err != nil {
		c.Error(errs.RenderError)
	}
}

func (c *Context) RenderBytes(s int, b []byte) {
	c.Status(s)

	_, err := c.WriteBytes(b)
	if err != nil {
		c.Error(errs.RenderError)
	}
}

func (c *Context) JSON(s int, body any) {
	writeContentType(c.writer, jsonContentType)
	j, err := json.Marshal(body)
	if err != nil {
		s = http.StatusBadRequest
	}

	c.RenderBytes(s, j)
}

func (c *Context) Param(k string) (string, error) {
	p, f := c.params[k]
	if !f {
		return "", errors.New("key not found")
	}
	return p, nil
}

func (c *Context) ParamInt(k string) (int, error) {
	p, f := c.params[k]
	if !f {
		return 0, errors.New("key not found")
	}

	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (c *Context) Query(k string) (string, error) {
	q := c.request.URL.Query()
	v := q[k]
	if len(v) < 1 {
		return "", errors.Join(errs.QueryError.Error, errors.New("unable to find"+k))
	}

	return v[0], nil
}
