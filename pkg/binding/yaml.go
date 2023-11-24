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

func (yamlBinding) Bind(r *http.Request, a any) error {
	return decodeYAML(r.Body, a)
}

func (yamlBinding) BindBody(b []byte, a any) error {
	return decodeYAML(bytes.NewReader(b), a)
}

func decodeYAML(r io.Reader, a any) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(a); err != nil {
		return err
	}
	return validate(a)
}
