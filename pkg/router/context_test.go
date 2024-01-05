package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var MockRoute = Route{
	Path:          "/mock/root",
	Params:        "/:mock",
	Method:        Get,
	HandlerFunc:   MockHandler(),
	DecoratorFunc: MockDecorator(),
}

func MockHandler() HandlerFunc {
	return func(c *Context) {
	}
}

func MockDecorator() DecoratorFunc {
	return func(f HandlerFunc) HandlerFunc {
		return func(c *Context) {
		}
	}
}

func TestContextReset(t *testing.T) {
	c := Context{}
	c.Reset()
}

func TestContextCopy(t *testing.T) {
	c := Context{}
	_ = c.Copy()
}

func TestContextSet(t *testing.T) {
	c := Context{}
	c.Set("key", "value")
	_, e := c.Get("key")
	assert.True(t, e)
}

func TestContextGet(t *testing.T) {
	c := Context{}
	c.Set("key", "value")
	v, _ := c.Get("key")
	assert.Equal(t, "value", v)
}

func TestContextMustGet(t *testing.T) {
	c := Context{}
	c.Set("key", "value")
	v, err := c.MustGet("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

	_, err = c.MustGet("keys")
	assert.Error(t, err)
}

func TestContextBody(t *testing.T) {

}

func TestContextMethod(t *testing.T) {
	c := Context{
		route: &MockRoute,
	}
	m := c.Method()
	assert.Equal(t, Get, m)
}

func TestContextError(t *testing.T) {

}

func TestContextErrorLog(t *testing.T) {

}

func TestContextWrite(t *testing.T) {

}

func TestContextWriteBytes(t *testing.T) {

}

func TestContextHeader(t *testing.T) {

}

func TestContextGetHeader(t *testing.T) {

}

func TestContextRoute(t *testing.T) {
	c := Context{
		route: &MockRoute,
	}

	r := c.Route()
	assert.Equal(t, MockRoute.rawPath, r.rawPath)
}

func TestContextPath(t *testing.T) {
	c := Context{
		route: &MockRoute,
	}

	p := c.Path()
	assert.Equal(t, "/mock/root/:mock", p)
}

func TestContextHandler(t *testing.T) {
	c := Context{
		route: &MockRoute,
	}

	f := c.Handler()
	assert.NotNil(t, f)
}

func TestContextDecorator(t *testing.T) {
	c := Context{
		route: &MockRoute,
	}

	f := c.Decorator()
	assert.NotNil(t, f)
}

func TestContextAborted(t *testing.T) {
	c := Context{}

	a := c.Aborted()
	assert.False(t, a)
}

func TestContextAbort(t *testing.T) {
	c := Context{}

	c.Abort()
	a := c.Aborted()
	assert.True(t, a)
}

func TestContextAbortWithStatus(t *testing.T) {
	// c := Context{}

	// c.AbortWithStatus(http.StatusOK)
	// a := c.Aborted()
	// assert.True(t, a)
}

func TestContextAbortWithError(t *testing.T) {
	// c := Context{}

	// c.AbortWithError(http.StatusBadRequest, errors.New("error"))
	// a := c.Aborted()
	// assert.True(t, a)
}

func TestContextStatus(t *testing.T) {

}

func TestContextRedirect(t *testing.T) {

}

func TestContextRender(t *testing.T) {

}

func TestContextRenderBytes(t *testing.T) {

}

func TestContextString(t *testing.T) {

}

func TestContextJSON(t *testing.T) {

}

func TestContextXML(t *testing.T) {

}

func TestContextYAML(t *testing.T) {

}

func TestContextTOML(t *testing.T) {

}

func TestContextHTML(t *testing.T) {

}

func TestContextGetFile(t *testing.T) {

}

func TestContextSaveFile(t *testing.T) {

}

func TestContextFormFile(t *testing.T) {

}

func TestContextMultipartForm(t *testing.T) {

}

func TestContextCookie(t *testing.T) {

}

func TestContextParam(t *testing.T) {

}

func TestContextParamInt(t *testing.T) {

}

func TestContextParamFloat(t *testing.T) {

}

func TestContextParamBool(t *testing.T) {

}

func TestContextParamAll(t *testing.T) {

}

func TestContextQuery(t *testing.T) {

}

func TestContextQueryInt(t *testing.T) {

}

func TestContextQueryFloat(t *testing.T) {

}

func TestContextQueryBool(t *testing.T) {

}

func TestContextQueryAll(t *testing.T) {

}

func TestContextRequestAddress(t *testing.T) {

}

func TestContextSecure(t *testing.T) {

}

func TestContextProtocol(t *testing.T) {

}

func TextContextShouldBindWith(t *testing.T) {

}

func TestContextMustBindWith(t *testing.T) {

}

func TextContextShouldBindJSON(t *testing.T) {

}

func TestContextBindJSON(t *testing.T) {

}

func TextContextShouldBindTOML(t *testing.T) {

}

func TestContextBindTOML(t *testing.T) {

}

func TextContextShouldBindXML(t *testing.T) {

}

func TestContextBindXML(t *testing.T) {

}

func TextContextShouldBindYAML(t *testing.T) {

}

func TestContextBindYAML(t *testing.T) {

}
