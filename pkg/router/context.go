package router

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joseph-beck/routey/pkg/binding"
	errs "github.com/joseph-beck/routey/pkg/error"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
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
	app   *App
	route *Route

	writer  http.ResponseWriter
	request *http.Request
	params  map[string]string
	state   State
	values  map[string]any
	mu      sync.Mutex

	queryCache  url.Values
	queryCached bool
}

// Reset the current Context
func (c *Context) Reset() {
	c.app = nil
	c.route = nil

	c.params = nil
	c.state = Healthy
	c.values = nil
	c.mu = sync.Mutex{}

	c.queryCache = nil
	c.queryCached = false
}

// Copy the current Context and give a pointer to the copy
func (c *Context) Copy() *Context {
	return &Context{
		app:   c.app,
		route: c.route,

		writer:  c.writer,
		request: c.request,
		params:  c.params,
		state:   c.state,
		values:  c.values,
		mu:      sync.Mutex{},

		queryCache:  c.queryCache,
		queryCached: c.queryCached,
	}
}

// Sets a value within the context
func (c *Context) Set(k string, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.values == nil {
		c.values = make(map[string]any)
	}

	c.values[k] = v
}

// Sets a value only if it does not exist
func (c *Context) MustSet(k string, v string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.values == nil {
		c.values = make(map[string]any)
	}

	_, e := c.values[k]
	if e {
		return errs.DataExistsError.Error
	}

	c.values[k] = v
	return nil
}

// Gets a value stored within the context
func (c *Context) Get(k string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, e := c.values[k]
	return v, e
}

// A value must be returned from this otherwise an error will occur
func (c *Context) MustGet(k string) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, e := c.values[k]
	if !e {
		return nil, errs.NoDataError.Error
	}

	return v, nil
}

// Get the body of the request
func (c *Context) Body() ([]byte, error) {
	if c.request.Body == nil {
		return nil, errors.New("empty body")
	}

	b, err := io.ReadAll(c.request.Body)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(b)
	c.request.Body = io.NopCloser(buff)
	n := bytes.NewBuffer(b)

	return n.Bytes(), nil
}

// Get the method of the route
func (c *Context) Method() Method {
	return c.route.Method
}

// Log an error to console
func (c *Context) Error(err errs.Error) {
	if err.Error == nil {
		panic("err is nil")
	}
}

// Log an error to console
func (c *Context) ErrorLog(err errs.Error) {
	if err.Error == nil {
		panic("err is nil")
	}

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

// Get the route of the matched route
func (c *Context) Route() *Route {
	if c.route == nil {
		return nil
	}

	return c.route.Copy()
}

// Get the path of the context
func (c *Context) Path() string {
	return c.route.Path + c.route.Params
}

// Get the Handler of the context
func (c *Context) Handler() HandlerFunc {
	return c.route.HandlerFunc
}

// Get the Decorator of the context
func (c *Context) Decorator() DecoratorFunc {
	return c.route.DecoratorFunc
}

// Has the context been aborted?
func (c *Context) Aborted() bool {
	return c.state == Aborted
}

// Aborts the current action
func (c *Context) Abort() {
	c.state = Aborted
}

// Abort with a status
func (c *Context) AbortWithStatus(s int) {
	c.state = Aborted
	c.Status(s)
}

// Abort with a status and an error
func (c *Context) AbortWithError(s int, e error) {
	c.state = Aborted
	c.Status(s)
}

// Respond with just a status
func (c *Context) Status(s int) {
	c.writer.WriteHeader(s)
}

// Redirects to the given location with the given status
func (c *Context) Redirect(s int, l string) {
	i := Redirect{
		status:   s,
		request:  c.request,
		location: l,
	}

	err := i.Render(c.writer)
	if err != nil {
		c.ErrorLog(errs.RedirectError)
		c.Abort()
	}
}

// Render a string body with status
func (c *Context) Render(s int, b string) {
	c.Status(s)

	_, err := c.Write(b)
	if err != nil {
		c.ErrorLog(errs.RenderError)
		c.Abort()
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

// Render a string with optional values
func (c *Context) String(s int, f string, v ...any) {
	writeContentType(c.writer, plainContentType)

	if len(v) > 0 {
		r := fmt.Sprintf(f, v)
		c.RenderBytes(s, []byte(r))
		return
	}

	c.RenderBytes(s, []byte(f))
}

// Render response in a JSON format from a body
func (c *Context) JSON(s int, b any) {
	writeContentType(c.writer, jsonContentType)

	j, err := json.Marshal(b)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.RenderBytes(s, j)
}

// Render XML
func (c *Context) XML(s int, b any) {
	writeContentType(c.writer, xmlContentType)

	x, err := xml.Marshal(b)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.RenderBytes(s, x)
}

// Render YAML
func (c *Context) YAML(s int, b any) {
	writeContentType(c.writer, yamlContentType)

	y, err := yaml.Marshal(b)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.RenderBytes(s, y)
}

// Render TOML
func (c *Context) TOML(s int, b any) {
	writeContentType(c.writer, tomlContentType)

	t, err := toml.Marshal(b)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.RenderBytes(s, t)
}

// Render HTML with a given file
func (c *Context) HTML(s int, n string, d any) {
	i := c.app.htmlRender.Instance(n, d)
	c.Status(s)

	err := i.Render(c.writer)
	if err != nil {
		c.Error(errs.HTMLError)
		c.Abort()
	}
}

// Serve the user a local file
func (c *Context) GetFile(p string) {
	http.ServeFile(c.writer, c.request, p)
}

// Save a file that was uploaded
func (c *Context) SaveFile(f *multipart.FileHeader, p string) error {
	s, err := f.Open()
	if err != nil {
		return err
	}
	defer s.Close()

	err = os.MkdirAll(filepath.Dir(p), 0750)
	if err != nil {
		return err
	}

	o, err := os.Create(p)
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = io.Copy(o, s)
	if err != nil {
		return err
	}

	return nil
}

// Get a form file
func (c *Context) FormFile(n string) (*multipart.FileHeader, error) {
	if c.request.MultipartForm == nil {
		err := c.request.ParseMultipartForm(32 << 30)
		if err != nil {
			return nil, err
		}
	}
	f, h, err := c.request.FormFile(n)
	if err != nil {
		return nil, err
	}

	f.Close()
	return h, nil
}

// Get a multipart form
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.request.ParseMultipartForm(32 << 30)
	if err != nil {
		return nil, err
	}
	f := c.request.MultipartForm

	return f, nil
}

// Get a value from the cookies of the request
func (c *Context) Cookie(n string) (string, error) {
	s, err := c.request.Cookie(n)
	if err != nil {
		return "", err
	}

	v, err := url.QueryUnescape(s.Value)
	if err != nil {
		return "", err
	}

	return v, nil
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

// Get a float value from the params
func (c *Context) ParamFloat(k string) (float64, error) {
	p, f := c.params[k]
	if !f {
		return 0, errors.New("key not found")
	}

	v, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

// Get a boolean value from the params
func (c *Context) ParamBool(k string) (bool, error) {
	p, f := c.params[k]
	if !f {
		return false, errors.New("key not found")
	}

	b, err := strconv.ParseBool(p)
	if err != nil {
		return false, err
	}

	return b, nil
}

// Get all of the params of the request
func (c *Context) ParamsAll() (map[string]string, error) {
	return c.params, nil
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

// Get a float value from the query
func (c *Context) QueryFloat(k string) (float64, error) {
	c.getQueryCache()
	v := c.queryCache[k]
	if len(v) < 1 {
		return 0, errors.Join(errs.QueryError.Error, errors.New("unable to find "+k))
	}

	f, err := strconv.ParseFloat(v[0], 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

// Get a boolean value from the query
func (c *Context) QueryBool(k string) (bool, error) {
	c.getQueryCache()
	v := c.queryCache[k]
	if len(v) < 1 {
		return false, errors.Join(errs.QueryError.Error, errors.New("unable to find "+k))
	}

	b, err := strconv.ParseBool(v[0])
	if err != nil {
		return false, err
	}

	return b, nil
}

// Get all the queries from the request
func (c *Context) QueryAll() (map[string]string, error) {
	c.getQueryCache()
	r := make(map[string]string)

	for k, v := range c.queryCache {
		if len(v) > 0 {
			r[k] = v[0]
		}
	}

	return r, nil
}

// Get the IP off the requester
func (c *Context) RequestAddress() (string, error) {
	h := c.request.Header.Get("X-Forwarded-For")
	if h != "" {
		a := strings.Split(h, ",")
		return strings.TrimSpace(a[0]), nil
	}

	a := strings.Split(c.request.RemoteAddr, ":")
	if a[0] == "" {
		return "", errs.HTMLError.Error
	}

	return a[0], nil
}

// Is the request secure?
func (c *Context) Secure() bool {
	return c.Protocol() == "https"
}

// Get the protocol used by the request
func (c *Context) Protocol() string {
	if c.request.TLS == nil {
		return "http"
	}

	return "https"
}

// Using the provided Binder, binds the context body with the a any
func (c *Context) ShouldBindWith(a any, b binding.Binder) error {
	return b.Bind(c.request, a)
}

// Must bind with the given struct
func (c *Context) MustBindWith(a any, b binding.Binder) error {
	err := c.ShouldBindWith(a, b)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return err
	}

	return nil
}

// Bind JSON to a any
func (c *Context) BindJSON(a any) error {
	return c.MustBindWith(a, binding.JSON)
}

// Wrapper for ShouldBindWith(a, binding.JSON)
func (c *Context) ShouldBindJSON(a any) error {
	return c.ShouldBindWith(a, binding.JSON)
}

// Bind TOML to a any
func (c *Context) BindTOML(a any) error {
	return c.MustBindWith(a, binding.TOML)
}

// Wrapper for ShouldBindWith(a, binding.TOML)
func (c *Context) ShouldBindTOML(a any) error {
	return c.ShouldBindWith(a, binding.TOML)
}

// Bind XML to a any
func (c *Context) BindXML(a any) error {
	return c.MustBindWith(a, binding.XML)
}

// Wrapper for ShouldBindWith(a, binding.XML)
func (c *Context) ShouldBindXML(a any) error {
	return c.ShouldBindWith(a, binding.XML)
}

// Bind YAML to a any
func (c *Context) BindYAML(a any) error {
	return c.MustBindWith(a, binding.YAML)
}

// Wrapper for ShouldBindWith(a, binding.YAML)
func (c *Context) ShouldBindYAML(a any) error {
	return c.ShouldBindWith(a, binding.YAML)
}
