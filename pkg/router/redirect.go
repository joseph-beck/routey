package router

import (
	"fmt"
	"net/http"
)

// Redirect struct
type Redirect struct {
	status   int
	request  *http.Request
	location string
}

// Render the redirect, just sends the requester to the specified location
func (r Redirect) Render(w http.ResponseWriter) error {
	if (r.status < http.StatusMultipleChoices || r.status > http.StatusPermanentRedirect) && r.status != http.StatusCreated {
		panic(fmt.Sprintf("Cannot redirect with status code %d", r.status))
	}

	http.Redirect(w, r.request, r.location, r.status)
	return nil
}

// Write the content type, there is no content
func (r Redirect) WriteContentType(http.ResponseWriter) {}
