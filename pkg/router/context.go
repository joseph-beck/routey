package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/joseph-beck/routey/pkg/binding"
	errs "github.com/joseph-beck/routey/pkg/error"
)

// A Context provides
//
//   - writer: write data back to the user
//
//   - request: request data that was given
//
//   - params: the parameters of the request
//
//   - queryCache: a cache of queries for this request
//
//   - queryCached: has the queryCache been made?
type Context struct {
	writer  http.ResponseWriter
	request *http.Request
	params  map[string]string

	queryCache  url.Values
	queryCached bool
}

// Log an error to console
func (c *Context) ErrorLog(err errs.Error) {
	log.Println(err.String())
}

// Write a string
func (c *Context) Write(s string) (int, error) {
	b, err := c.writer.Write([]byte(s))
	return b, err
}

// Write bytes
func (c *Context) WriteBytes(body []byte) (int, error) {
	b, err := c.writer.Write(body)
	return b, err
}

// Add to the header of a response with a key and value
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

// Get the header of a request using a key
func (c *Context) GetHeader(k string) string {
	v := c.request.Header.Get(k)
	return v
}

// Respond with just a status
func (c *Context) Status(s int) {
	c.writer.WriteHeader(s)
}

// Render a string body with status
func (c *Context) Render(s int, b string) {
	c.Status(s)

	_, err := c.Write(b)
	if err != nil {
		c.ErrorLog(errs.RenderError)
	}
}

// Render a byte array with status
func (c *Context) RenderBytes(s int, b []byte) {
	c.Status(s)

	_, err := c.WriteBytes(b)
	if err != nil {
		c.ErrorLog(errs.RenderError)
	}
}

// Render response in a JSON format from a body
func (c *Context) JSON(s int, b any) {
	writeContentType(c.writer, jsonContentType)
	j, err := json.Marshal(b)
	if err != nil {
		s = http.StatusBadRequest
	}

	c.RenderBytes(s, j)
}

// Get a string from the parameters using a key
func (c *Context) Param(k string) (string, error) {
	p, f := c.params[k]
	if !f {
		return "", errors.New("key not found")
	}
	return p, nil
}

// Get an integer from the parameters using a key
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

// Get the query cache
func (c *Context) getQueryCache() {
	if c.queryCached {
		return
	}

	c.queryCache = c.request.URL.Query()
	c.queryCached = true
}

// Query a string from the context's query
func (c *Context) Query(k string) (string, error) {
	c.getQueryCache()
	v := c.queryCache[k]
	if len(v) < 1 {
		return "", errors.Join(errs.QueryError.Error, errors.New("unable to find "+k))
	}

	return v[0], nil
}

// Query an integer from the context's query
func (c *Context) QueryInt(k string) (int, error) {
	c.getQueryCache()
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

// Using the provided Binder, binds the context body with the a any
func (c *Context) ShouldBindWith(a any, b binding.Binder) error {
	return b.Bind(c.request, a)
}

// Bind JSON to a any
func (c *Context) BindJSON(a any) error {
	return nil
}

// Wrapper for ShouldBindWith(a, binding.JSON)
func (c *Context) ShouldBindJSON(a any) error {
	return c.ShouldBindWith(a, binding.JSON)
}

// Bind TOML to a any
func (c *Context) BindTOML(a any) error {
	return nil
}

// Wrapper for ShouldBindWith(a, binding.TOML)
func (c *Context) ShouldBindTOML(a any) error {
	return c.ShouldBindWith(a, binding.TOML)
}

// Bind XML to a any
func (c *Context) BindXML(a any) error {
	return nil
}

// Wrapper for ShouldBindWith(a, binding.XML)
func (c *Context) ShouldBindXML(a any) error {
	return c.ShouldBindWith(a, binding.XML)
}

// Bind YAML to a any
func (c *Context) BindYAML(a any) error {
	return nil
}

// Wrapper for ShouldBindWith(a, binding.YAML)
func (c *Context) ShouldBindYAML(a any) error {
	return c.ShouldBindWith(a, binding.YAML)
}
