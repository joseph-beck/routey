package error

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorNew(t *testing.T) {
	n := "Error"
	e := errors.New("error")
	err := New(n, e)

	assert.NotNil(t, err)
}

func TestErrorString(t *testing.T) {
	n := "Error"
	e := errors.New("error")
	err := New(n, e)

	assert.Equal(t, n, err.String())
}

func TestErrorEqual(t *testing.T) {
	assert.True(t, RenderError.Equal(RenderError))
}

func TestErrorNotEqual(t *testing.T) {
	assert.False(t, RenderError.Equal(HTMLError))
}
