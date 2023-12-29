package binding

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/pelletier/go-toml/v2"
)

type tomlBinding struct{}

func (t tomlBinding) Name() string {
	return "toml"
}

func (t tomlBinding) Bind(r *http.Request, a any) error {
	buff, err := t.readBody(r)
	if err != nil {
		return err
	}

	return t.decodeToml(buff, a)
}

func (t tomlBinding) BindBody(b []byte, a any) error {
	return t.decodeToml(bytes.NewReader(b), a)
}

func (t tomlBinding) decodeToml(r io.Reader, a any) error {
	decoder := toml.NewDecoder(r)

	err := decoder.Decode(a)
	if err != nil {
		return err
	}

	return decoder.Decode(a)
}

func (t tomlBinding) readBody(r *http.Request) (*bytes.Buffer, error) {
	if r.Body == nil {
		return nil, errors.New("invalid request")
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(b)
	r.Body = io.NopCloser(buff)
	n := bytes.NewBuffer(b)
	return n, nil
}
