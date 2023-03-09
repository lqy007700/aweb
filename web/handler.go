package web

type HandlerFunc func(c *Context)

type Handler interface {
	ServeHTTP(ctx *Context)
	Routable
}
