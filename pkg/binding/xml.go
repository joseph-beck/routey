package binding

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

type xmlBinding struct{}

func (x xmlBinding) Name() string {
	return "xml"
}

func (x xmlBinding) Bind(r *http.Request, a any) error {
	buff, err := x.readBody(r)
	if err != nil {
		return err
	}

	return x.decodeXML(buff, a)
}

func (x xmlBinding) BindBody(b []byte, a any) error {
	return x.decodeXML(bytes.NewReader(b), a)
}
func (x xmlBinding) decodeXML(r io.Reader, a any) error {
	decoder := xml.NewDecoder(r)

	err := decoder.Decode(a)
	if err != nil {
		return err
	}

	return validate(a)
}

func (x xmlBinding) readBody(r *http.Request) (*bytes.Buffer, error) {
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
