package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	errs "github.com/joseph-beck/routey/pkg/error"
)

type Context struct {
	writer  http.ResponseWriter
	request *http.Request
	params  map[string]string

	queryCache  url.Values
	queryCached bool
}

func (c *Context) ErrorLog(err errs.Error) {
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
	if k == "" {
		return
	}

	if v == "" {
		c.writer.Header().Del(k)
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
		c.ErrorLog(errs.RenderError)
	}
}

func (c *Context) RenderBytes(s int, b []byte) {
	c.Status(s)

	_, err := c.WriteBytes(b)
	if err != nil {
		c.ErrorLog(errs.RenderError)
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

func (c *Context) queryCacheGet() {
	if c.queryCached {
		return
	}

	c.queryCache = c.request.URL.Query()
	c.queryCached = true
}

func (c *Context) Query(k string) (string, error) {
	c.queryCacheGet()
	v := c.queryCache[k]
	if len(v) < 1 {
		return "", errors.Join(errs.QueryError.Error, errors.New("unable to find "+k))
	}

	return v[0], nil
}

func (c *Context) QueryInt(k string) (int, error) {
	c.queryCacheGet()
	v := c.queryCache[k]
	if len(v) < 1 {
		return 0, errors.Join(errs.QueryError.Error, errors.New("unable to find "+k))
	}

	i, err := strconv.Atoi(v[0])
	if err != nil {
		return 0, err
	}

	return i, nil
}
