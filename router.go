package main

type Routable interface {
	Route(method, pattern string, handler HandlerFunc) error
}
