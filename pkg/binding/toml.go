package binding

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pelletier/go-toml/v2"
)

type tomlBinding struct{}

func (t tomlBinding) Name() string {
	return "toml"
}

func (t tomlBinding) Bind(r *http.Request, a any) error {
	return t.decodeToml(r.Body, a)
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
