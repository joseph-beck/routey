package binding

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonBinding struct{}

func (j jsonBinding) Name() string {
	return "json"
}

func (j jsonBinding) Bind(r *http.Request, a any) error {
	buff, err := j.readBody(r)
	if err != nil {
		return err
	}

	return j.decodeJSON(buff, a)
}

func (j jsonBinding) BindBody(b []byte, a any) error {
	return j.decodeJSON(bytes.NewReader(b), a)
}

func (j jsonBinding) decodeJSON(r io.Reader, a any) error {
	decoder := json.NewDecoder(r)

	err := decoder.Decode(a)
	if err != nil {
		return err
	}
	return validate(a)
}

func (j jsonBinding) readBody(r *http.Request) (*bytes.Buffer, error) {
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
