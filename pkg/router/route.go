package router

import (
	"regexp"
)

type Route struct {
	Path          string
	Params        string
	Method        Method
	HandlerFunc   HandlerFunc
	DecoratorFunc DecoratorFunc

	regexp *regexp.Regexp
}

func (route *Route) Match(c *Context) bool {
	method := parseMethod(c.request.Method)
	if method != route.Method {
		return false
	}

	route.regexp = regexp.MustCompile("^" + route.Path + route.Params + "$")
	match := route.regexp.FindStringSubmatch(c.request.URL.Path)
	if match == nil {
		return false
	}

	params := make(map[string]string)
	groups := route.regexp.SubexpNames()
	for i, group := range match {
		params[groups[i]] = group
	}
	c.params = params

	return true
}
