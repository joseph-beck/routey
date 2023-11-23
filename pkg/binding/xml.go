package binding

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type xmlBinding struct{}

func (x xmlBinding) Name() string {
	return "xml"
}

func (x xmlBinding) Bind(r *http.Request, a any) error {
	return x.decodeXML(r.Body, a)
}

func (x xmlBinding) BindBody(b []byte, a any) error {
	return x.decodeXML(bytes.NewReader(b), a)
}
func (x xmlBinding) decodeXML(r io.Reader, a any) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(a); err != nil {
		return err
	}
	return validate(a)
}
