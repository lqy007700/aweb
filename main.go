package main

import (
	"net/http"
)

type Server interface {
	Route(method, pattern string, handle func(ctx *Context))
	Start(addr string)
}

type aWebServer struct {
	name     string
	handlers *HandleBaseOnMap
}

func newAWebServer() Server {
	return &aWebServer{
		name:     "app",
		handlers: NewHandleBaseOnMap(),
	}
}

func (a *aWebServer) Route(method, pattern string, handle func(ctx *Context)) {
	key := a.handlers.key(method, pattern)
	a.handlers.router[key] = handle
}

func (a *aWebServer) Start(addr string) {
	http.Handle("/", a.handlers)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	srv := newAWebServer()
	srv.Route("POST", "/sign", signUp)
	srv.Start(":8080")
}
