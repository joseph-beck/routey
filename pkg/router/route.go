package router

import "net/http"

type Route struct {
	Path          string
	Params        string
	Method        Method
	HandlerFunc   HandlerFunc
	DecoratorFunc DecoratorFunc
}

func parseMethod(s string) Method {
	switch s {
	case "GET":
		return Get
	case "POST":
		return Post
	case "PUT":
		return Put
	case "PATCH":
		return Patch
	case "DELEET":
		return Delete
	case "HEAD":
		return Head
	case "OPTIONS":
		return Options
	default:
		return Undefined
	}
}

func (route *Route) Match(r *http.Request) bool {
	m := parseMethod(r.Method)
	if m != route.Method {
		return false
	}

	if r.URL.Path != string(route.Path)+string(route.Params) {
		return false
	}

	return true
}
