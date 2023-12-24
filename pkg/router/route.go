package router

import (
	"regexp"
)

// Route struct
//
//   - Path: path of the route request
//
//   - Params: params of the route start with :
//
//   - Method: method of the route
//
//   - HandlerFunc: handler function of the route
//
//   - DecoratorFunc: decorator function of the route
//
//   - regexp: regexp used for params
type Route struct {
	Path          string
	Params        string
	Method        Method
	HandlerFunc   HandlerFunc
	DecoratorFunc DecoratorFunc

	regexp    *regexp.Regexp
	rawPath   string
	formatted bool
}

// Match a Route with a Context
func (r *Route) Match(c *Context) bool {
	method := parseMethod(c.request.Method)
	if method != r.Method {
		return false
	}

	r.regexp = regexp.MustCompile("^" + r.rawPath + "$")
	match := r.regexp.FindStringSubmatch(c.request.URL.Path)
	if match == nil {
		return false
	}

	params := make(map[string]string)
	groups := r.regexp.SubexpNames()
	for i, group := range match {
		params[groups[i]] = group
	}
	c.params = params

	return true
}

// Format the params of the route
func (r *Route) Format() error {
	if r.formatted {
		return nil
	}

	p, err := parsePathParams(r.Path + r.Params)
	if err != nil {
		return err
	}
	r.rawPath = p

	return nil
}
