package websockets

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/joseph-beck/routey/pkg/binding"
	routey "github.com/joseph-beck/routey/pkg/router"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v3"
)

func defaultRecover(c *WSConn) {
	if err := recover(); err != nil {
		_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", err, debug.Stack())) //nolint:errcheck // This will never fail

		err := c.SendJSON(routey.M{"error": err})
		if err != nil {
			_, _ = os.Stderr.WriteString(fmt.Sprintf("could not write error response: %v\n", err)) //nolint:errcheck // This will never fail
		}
	}
}

type WSConn struct {
	socket  *websocket.Conn
	params  map[string]string
	cookies map[string]string
	headers map[string]string
	queries map[string]string
	ip      string
}

func New(c ...Config) *WSConn {
	return &WSConn{}
}

// Get a param from the ws connection
func (c *WSConn) Params(k string, d ...string) (string, error) {
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
func (c *WSConn) ParamsInt(k string, d ...int) (int, error) {
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
func (c *WSConn) ParamsFloat(k string, d ...float64) (float64, error) {
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

func (c *WSConn) ParamsBool(k string, d ...bool) (bool, error) {
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

func (c *WSConn) Query(k string, d ...string) (string, error) {
	v, ok := c.queries[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return "", errors.New("error, no params found, maybe give it a default?")
	}

	return v, nil
}

func (c *WSConn) QueryInt(k string, d ...int) (int, error) {
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

func (c *WSConn) QueryFloat(k string, d ...float64) (float64, error) {
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

func (c *WSConn) QueryBool(k string, d ...bool) (bool, error) {
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

func (c *WSConn) Cookies(k string, d ...string) (string, error) {
	v, ok := c.cookies[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return "", errors.New("error, no params found, maybe give it a default?")
	}

	return v, nil
}

func (c *WSConn) Headers(k string, d ...string) (string, error) {
	v, ok := c.headers[k]
	if !ok && len(d) > 0 {
		return d[0], nil
	}

	if !ok && len(d) <= 0 {
		return "", errors.New("error, no params found, maybe give it a default?")
	}

	return v, nil
}

func (c *WSConn) IP() string {
	return c.ip
}

func (c *WSConn) Read(b []byte) error {
	_, err := c.socket.Read(b)
	if err != nil {
		return err
	}

	return nil
}

func (c *WSConn) ReadMessage(m string) error {
	var d []byte
	_, err := c.socket.Read(d)
	if err != nil {
		return err
	}

	return nil
}

func (c *WSConn) ReadJSON(b any) error {
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

func (c *WSConn) ReadXML(b any) error {
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

func (c *WSConn) ReadYAML(b any) error {
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

func (c *WSConn) ReadTOML(b any) error {
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

func (c *WSConn) Send(b []byte) error {
	_, err := c.socket.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (c *WSConn) SendMessage(m string) error {
	b := []byte(m)
	_, err := c.socket.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (c *WSConn) SendJSON(b any) error {
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

func (c *WSConn) SendXML(b any) error {
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

func (c *WSConn) SendYAML(b any) error {
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

func (c *WSConn) SendTOML(b any) error {
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

func (c *WSConn) Close() error {
	err := c.socket.Close()
	if err != nil {
		return err
	}

	return nil
}
