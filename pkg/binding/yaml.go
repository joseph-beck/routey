package binding

import (
	"bytes"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (y yamlBinding) Bind(r *http.Request, a any) error {
	return y.decodeYAML(r.Body, a)
}

func (y yamlBinding) BindBody(b []byte, a any) error {
	return y.decodeYAML(bytes.NewReader(b), a)
}

func (y yamlBinding) decodeYAML(r io.Reader, a any) error {
	decoder := yaml.NewDecoder(r)

	err := decoder.Decode(a)
	if err != nil {
		return err
	}

	return validate(a)
}
