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
	if r == nil || r.Body == nil {
		return errors.New("invalid request")
	}

	return j.decodeJSON(r.Body, a)
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
