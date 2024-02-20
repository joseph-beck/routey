package websockets

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/joseph-beck/routey/pkg/binding"
	routey "github.com/joseph-beck/routey/pkg/router"
	"golang.org/x/net/websocket"
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
	return "", nil
}

func (c *WSConn) Headers(k string, d ...string) (string, error) {
	return "", nil
}

func (c *WSConn) IP() string {
	return c.ip
}

func (c *WSConn) Read(b []byte) error {
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

func (c *WSConn) ReadJSON(i interface{}) error {
	var d []byte
	_, err := c.socket.Read(d)
	if err != nil {
		return err
	}

	err = binding.JSON.BindBody(d, i)
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
func (c *WSConn) SendJSON(i interface{}) error {
	return nil
}

func (c *WSConn) Close() error {
	err := c.socket.Close()
	if err != nil {
		return err
	}

	return nil
}
