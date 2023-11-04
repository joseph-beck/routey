package router

type HandlerFunc func(c *Context)

func (f HandlerFunc) Serve(c *Context) {
	f(c)
}

type DecoratorFunc func(f HandlerFunc) HandlerFunc
