package router

import (
	"fmt"
	"net/http"
)

type Redirect struct {
	status   int
	request  *http.Request
	location string
}

func (r Redirect) Render(w http.ResponseWriter) error {
	if (r.status < http.StatusMultipleChoices || r.status > http.StatusPermanentRedirect) && r.status != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.status))
	}

	http.Redirect(w, r.request, r.location, r.status)
	return nil
}

func (r Redirect) WriteContentType(http.ResponseWriter) {}
