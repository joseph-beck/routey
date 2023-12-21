package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	m := Get
	s := m.String()

	assert.Equal(t, "GET", s)
}

func TestParseMethod(t *testing.T) {
	s := "GET"
	m := parseMethod(s)

	assert.Equal(t, Get, m)
}
