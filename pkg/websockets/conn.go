package websockets

import (
	"fmt"
	"os"
	"runtime/debug"

	routey "github.com/joseph-beck/routey/pkg/router"
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
	params  map[string]string
	cookies map[string]string
	headers map[string]string
	queries map[string]string
	ip      string
}

func (c *WSConn) Params(k string, d ...string) (string, error) {
	return d[0], nil
}

func (c *WSConn) ParamsInt(k string, d ...int) (int, error) {
	return d[0], nil
}

func (c *WSConn) ParamsFloat(k string, d ...float64) (float64, error) {
	return d[0], nil
}

func (c *WSConn) ParamsBool(k string, d ...bool) (bool, error) {
	return d[0], nil
}

func (c *WSConn) Query(k string, d ...string) (string, error) {
	return d[0], nil
}

func (c *WSConn) QueryInt(k string, d ...int) (int, error) {
	return d[0], nil
}

func (c *WSConn) QueryFloat(k string, d ...float64) (float64, error) {
	return d[0], nil
}

func (c *WSConn) QueryBool(k string, d ...bool) (bool, error) {
	return d[0], nil
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

func (c *WSConn) ReadMessage(m string) error {
	return nil
}

func (c *WSConn) SendMessage(m string) error {
	return nil
}

func (c *WSConn) ReadJSON(i interface{}) error {
	return nil
}

func (c *WSConn) SendJSON(i interface{}) error {
	return nil
}
