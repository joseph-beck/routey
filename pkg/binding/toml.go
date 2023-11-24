package binding

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pelletier/go-toml/v2"
)

type tomlBinding struct{}

func (tomlBinding) Name() string {
	return "toml"
}

func (tomlBinding) Bind(r *http.Request, a any) error {
	return decodeToml(r.Body, a)
}

func (tomlBinding) BindBody(b []byte, a any) error {
	return decodeToml(bytes.NewReader(b), a)
}

func decodeToml(r io.Reader, a any) error {
	decoder := toml.NewDecoder(r)
	if err := decoder.Decode(a); err != nil {
		return err
	}
	return decoder.Decode(a)
}
