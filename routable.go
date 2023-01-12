package main

type Routable interface {
	Route(method, pattern string, handle handleFunc)
}