package websockets

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/joseph-beck/routey/pkg/binding"
	routey "github.com/joseph-beck/routey/pkg/router"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v3"
)

// Default recovery function, just returns the stack trace
func defaultRecover(c *Conn) {
	if err := recover(); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", err, debug.Stack())) //nolint:errcheck // This will never fail

		err := c.SendJSON(routey.M{"error": err})
		if err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("could not write error response: %v\n", err)) //nolint:errcheck // This will never fail
		}
	}
}

// WebSocket connection
type Conn struct {
	socket  *websocket.Conn
	params  map[string]string
	cookies map[string]string
	headers map[string]string
	queries map[string]string
	ip      string
}

// Creates a new handler function of a given handler function for a WebSocket handler
func New(h func(*Conn), c ...Config) routey.HandlerFunc {
	cfg := Default()
	if len(c) > 0 {
		cfg = c[0]
	}

	if len(cfg.Origins) == 0 {
		cfg.Origins = []string{}
	}

	if cfg.ReadBufferSize == 0 {
		cfg.ReadBufferSize = 1024
	}

	if cfg.WriteBufferSize == 0 {
		cfg.WriteBufferSize = 1024
	}

	if cfg.RecoverHandler == nil {
		cfg.RecoverHandler = defaultRecover
	}

	return func(c *routey.Context) {
		c.Status(http.StatusOK)
	}
}

// Get a param from the ws connection
func (c *Conn) Params(k string, d ...string) (string, error) {
	v, ok := c.params[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return "", errors.New("error, no params found, maybe give it a default?")
	}

	return v, nil
}

// Get a param of type int from the ws connection
func (c *Conn) ParamsInt(k string, d ...int) (int, error) {
	v, ok := c.params[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return 0, errors.New("error, no params found, maybe give it a default?")
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// Get a param of type float from the ws connection
func (c *Conn) ParamsFloat(k string, d ...float64) (float64, error) {
	v, ok := c.params[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return 0, errors.New("error, no params found, maybe give it a default?")
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

// Get a param of type bool from the ws connection
func (c *Conn) ParamsBool(k string, d ...bool) (bool, error) {
	v, ok := c.params[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return false, errors.New("error, no params found, maybe give it a default?")
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, err
	}

	return b, nil
}

// Get a query from the ws connection
func (c *Conn) Query(k string, d ...string) (string, error) {
	v, ok := c.queries[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return "", errors.New("error, no params found, maybe give it a default?")
	}

	return v, nil
}

// Get a query of type int from the ws connection
func (c *Conn) QueryInt(k string, d ...int) (int, error) {
	v, ok := c.queries[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return 0, errors.New("error, no params found, maybe give it a default?")
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return i, nil
}

// Get a query of type float from the ws connection
func (c *Conn) QueryFloat(k string, d ...float64) (float64, error) {
	v, ok := c.queries[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return 0, errors.New("error, no params found, maybe give it a default?")
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

// Get a query of type bool from the ws connection
func (c *Conn) QueryBool(k string, d ...bool) (bool, error) {
	v, ok := c.queries[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return false, errors.New("error, no params found, maybe give it a default?")
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, err
	}

	return b, nil
}

// Get the value of a cookie from the ws connection
func (c *Conn) Cookies(k string, d ...string) (string, error) {
	v, ok := c.cookies[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return "", errors.New("error, no params found, maybe give it a default?")
	}

	return v, nil
}

// Get the value of a header from the ws connection
func (c *Conn) Headers(k string, d ...string) (string, error) {
	v, ok := c.headers[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return "", errors.New("error, no params found, maybe give it a default?")
	}

	return v, nil
}

// Get the ip of the connection from the ws
func (c *Conn) IP() string {
	return c.ip
}

// Read to a byte array from the
func (c *Conn) Read(b []byte) error {
	_, err := c.socket.Read(b)
	if err != nil {
		return err
	}

	return nil
}

// Read a message into a reference
func (c *Conn) ReadMessage(m string) error {
	var d []byte
	_, err := c.socket.Read(d)
	if err != nil {
		return err
	}

	m = string(d)

	return nil
}

// Read the value of the body into a JSON object
func (c *Conn) ReadJSON(b any) error {
	var d []byte
	_, err := c.socket.Read(d)
	if err != nil {
		return err
	}

	err = binding.JSON.BindBody(d, b)
	if err != nil {
		return err
	}

	return nil
}

// Read the value of the body into an XML object
func (c *Conn) ReadXML(b any) error {
	var d []byte
	_, err := c.socket.Read(d)
	if err != nil {
		return err
	}

	err = binding.XML.BindBody(d, b)
	if err != nil {
		return err
	}

	return nil
}

// Read the value of the body into a YAML object
func (c *Conn) ReadYAML(b any) error {
	var d []byte
	_, err := c.socket.Read(d)
	if err != nil {
		return err
	}

	err = binding.YAML.BindBody(d, b)
	if err != nil {
		return err
	}

	return nil
}

// Read the value of the body into a TOML object
func (c *Conn) ReadTOML(b any) error {
	var d []byte
	_, err := c.socket.Read(d)
	if err != nil {
		return err
	}

	err = binding.TOML.BindBody(d, b)
	if err != nil {
		return err
	}

	return nil
}

// Send a byte array through the ws
func (c *Conn) Send(b []byte) error {
	_, err := c.socket.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// Send a string message through the ws
func (c *Conn) SendMessage(m string) error {
	b := []byte(m)
	_, err := c.socket.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// Send a JSON object through the ws
func (c *Conn) SendJSON(b any) error {
	j, err := json.Marshal(b)
	if err != nil {
		return err
	}

	_, err = c.socket.Write(j)
	if err != nil {
		return err
	}

	return nil
}

// Send an XML object through the ws
func (c *Conn) SendXML(b any) error {
	j, err := xml.Marshal(b)
	if err != nil {
		return err
	}

	_, err = c.socket.Write(j)
	if err != nil {
		return err
	}

	return nil
}

// Send a YAML object through the ws
func (c *Conn) SendYAML(b any) error {
	j, err := yaml.Marshal(b)
	if err != nil {
		return err
	}

	_, err = c.socket.Write(j)
	if err != nil {
		return err
	}

	return nil
}

// Send a TOML object through the ws
func (c *Conn) SendTOML(b any) error {
	j, err := toml.Marshal(b)
	if err != nil {
		return err
	}

	_, err = c.socket.Write(j)
	if err != nil {
		return err
	}

	return nil
}

// Close the ws connection
func (c *Conn) Close() error {
	err := c.socket.Close()
	if err != nil {
		return err
	}

	return nil
}
