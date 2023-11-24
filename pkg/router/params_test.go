package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParamsOne(t *testing.T) {
	r, err := parseParams("/:param")
	assert.NoError(t, err)
	assert.Equal(t, r, "/(?P<param>\\w+)")

	r, err = parseParams("/:one/:two")
	assert.NoError(t, err)
	assert.Equal(t, r, "/(?P<one>\\w+)/(?P<two>\\w+)")
}

func TestParseParamsTwo(t *testing.T) {
	r, err := parseParams("")
	assert.NoError(t, err)
	assert.Equal(t, r, "")
}

func TestParsePathParamsOne(t *testing.T) {
	r, err := parsePathParams("/test/:param")
	assert.NoError(t, err)
	assert.Equal(t, r, "/test/(?P<param>\\w+)")

	r, err = parsePathParams("/test/:one/:two")
	assert.NoError(t, err)
	assert.Equal(t, r, "/test/(?P<one>\\w+)/(?P<two>\\w+)")
}

func TestParsePathParamsTwo(t *testing.T) {
	r, err := parsePathParams("")
	assert.NoError(t, err)
	assert.Equal(t, r, "/")
}
