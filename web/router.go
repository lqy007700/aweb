package web

type Routable interface {
	Route(method, pattern string, handler HandlerFunc, ms ...Middleware)
}
