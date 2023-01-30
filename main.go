package main

import (
	"net/http"
)

type Server interface {
	Routable
	Start(addr string)
}

type aWebServer struct {
	name    string
	handler Handle
	root    Filter
}

func newAWebServer(builders ...FilterBuilder) Server {
	handler := NewHandlerBaseOnTree()

	var root Filter = handler.ServeHTTP
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	return &aWebServer{
		name:    "app",
		handler: handler,
		root:    root,
	}
}

func (a *aWebServer) Route(method, pattern string, handle handleFunc) error {
	err := a.handler.Route(method, pattern, handle)
	return err
}

func (a *aWebServer) Start(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := newContext(w, r)
		a.root(c)
	})
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	srv := newAWebServer(Builder, Aprint)
	srv.Route("POST", "/sign/*", signUp)
	srv.Route("POST", "/sign/*/name", signUp1)

	srv.Start(":8080")
}
