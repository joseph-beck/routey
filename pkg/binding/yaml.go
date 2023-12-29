package binding

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (y yamlBinding) Bind(r *http.Request, a any) error {
	buff, err := y.readBody(r)
	if err != nil {
		return err
	}

	return y.decodeYAML(buff, a)
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

func (y yamlBinding) readBody(r *http.Request) (*bytes.Buffer, error) {
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
